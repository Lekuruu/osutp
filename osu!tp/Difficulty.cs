﻿using System;
using System.Collections.Generic;
using System.Linq;
using osu.GameplayElements.Beatmaps;
using osu.GameplayElements.HitObjects;

namespace osutp.TomPoints
{
    /// <summary>
    /// osu!tp's difficulty calculator ported to the osu! sdk as far as so far possible.
    /// </summary>
    public class TpDifficulty
    {
        // Those values are used as array indices. Be careful when changing them!
        public enum DifficultyType : int
        {
            Speed = 0,
            Aim,
        };

        private const double STAR_SCALING_FACTOR = 0.045;
        private const double EXTREME_SCALING_FACTOR = 0.5;
        private const float PLAYFIELD_WIDTH = 512;

        public TpDifficultyCalculation Process(BeatmapBase beatmap, List<HitObjectBase> hitObjects)
        {
            // Fill our custom tpHitObject class, that carries additional information
            List<TpHitObject> tpHitObjects = new List<TpHitObject>(hitObjects.Count);
            float CircleRadius = (PLAYFIELD_WIDTH / 16.0f) * (1.0f - 0.7f * ((float)beatmap.DifficultyCircleSize - 5.0f) / 5.0f);

            foreach (HitObjectBase hitObject in hitObjects)
            {
                tpHitObjects.Add(new TpHitObject(hitObject, CircleRadius));
            }

            // Sort tpHitObjects by StartTime of the HitObjects - just to make sure. Not using CompareTo, since it results in a crash (HitObjectBase inherits MarshalByRefObject)
            tpHitObjects.Sort((a, b) => a.BaseHitObject.StartTime - b.BaseHitObject.StartTime);

            if (!CalculateStrainValues(tpHitObjects))
                return null;

            double SpeedDifficulty = CalculateDifficulty(tpHitObjects, DifficultyType.Speed);
            double AimDifficulty = CalculateDifficulty(tpHitObjects, DifficultyType.Aim);

            // OverallDifficulty is not considered in this algorithm and neither is HpDrainRate. That means, that in this form the algorithm determines how hard it physically is
            // to play the map, assuming, that too much of an error will not lead to a death.
            // It might be desirable to include OverallDifficulty into map difficulty, but in my personal opinion it belongs more to the weighting of the actual peformance
            // and is superfluous in the beatmap difficulty rating.
            // If it were to be considered, then I would look at the hit window of normal HitCircles only, since Sliders and Spinners are (almost) "free" 300s and take map length
            // into account as well.

            // The difficulty can be scaled by any desired metric.
            // In osu!tp it gets squared to account for the rapid increase in difficulty as the limit of a human is approached. (Of course it also gets scaled afterwards.)
            // It would not be suitable for a star rating, therefore:

            // The following is a proposal to forge a star rating from 0 to 5. It consists of taking the square root of the difficulty, since by simply scaling the easier
            // 5-star maps would end up with one star.
            double SpeedStars = Math.Sqrt(SpeedDifficulty) * STAR_SCALING_FACTOR;
            double AimStars = Math.Sqrt(AimDifficulty) * STAR_SCALING_FACTOR;

            // Again, from own observations and from the general opinion of the community a map with high speed and low aim (or vice versa) difficulty is harder,
            // than a map with mediocre difficulty in both. Therefore we can not just add both difficulties together, but will introduce a scaling that favors extremes.
            double StarRating = SpeedStars + AimStars + Math.Abs(SpeedStars - AimStars) * EXTREME_SCALING_FACTOR;

            // Another approach to this would be taking Speed and Aim separately to a chosen power, which again would be equivalent. This would be more convenient if
            // the hit window size is to be considered as well.

            // Note: The star rating is tuned extremely tight! Airman (/b/104229) and Freedom Dive (/b/126645), two of the hardest ranked maps, both score ~4.66 stars.
            // Expect the easier kind of maps that officially get 5 stars to obtain around 2 by this metric. The tutorial still scores about half a star.
            // Tune by yourself as you please. ;)
            return new TpDifficultyCalculation
            {
                AmountNormal = tpHitObjects.FindAll(x => x.BaseHitObject.Type.HasFlag(HitObjectType.Normal)).Count,
                AmountSliders = tpHitObjects.FindAll(x => x.BaseHitObject.Type.HasFlag(HitObjectType.Slider)).Count,
                AmountSpinners = tpHitObjects.FindAll(x => x.BaseHitObject.Type.HasFlag(HitObjectType.Spinner)).Count,
                MaxCombo = tpHitObjects.Sum(x => x.Combo),
                SpeedDifficulty = SpeedDifficulty,
                AimDifficulty = AimDifficulty,
                StarRating = StarRating,
                SpeedStars = SpeedStars,
                AimStars = AimStars,
                OverallDifficulty = beatmap.DifficultyOverall,
                ApproachRate = beatmap.DifficultyApproachRate,
                HpDrainRate = beatmap.DifficultyHpDrainRate,
                CircleSize = beatmap.DifficultyCircleSize,
                SliderTickRate = beatmap.DifficultySliderTickRate,
                SliderMultiplier = beatmap.DifficultySliderMultiplier,
            };
        }
        
        // Exceptions would be nicer to handle errors, but for this small project it shall be ignored.
        private bool CalculateStrainValues(List<TpHitObject> tpHitObjects)
        {
            // Traverse hitObjects in pairs to calculate the strain value of NextHitObject from the strain value of CurrentHitObject and environment.
            List<TpHitObject>.Enumerator HitObjectsEnumerator = tpHitObjects.GetEnumerator();

            if (HitObjectsEnumerator.MoveNext() == false)
            {
                // Can not compute difficulty of empty beatmap
                return false;
            }

            TpHitObject CurrentHitObject = HitObjectsEnumerator.Current;
            TpHitObject NextHitObject;

            // First hitObject starts at strain 1. 1 is the default for strain values, so we don't need to set it here. See tpHitObject.
            while (HitObjectsEnumerator.MoveNext())
            {
                NextHitObject = HitObjectsEnumerator.Current;
                NextHitObject.CalculateStrains(CurrentHitObject);
                CurrentHitObject = NextHitObject;
            }

            return true;
        }

        // In milliseconds. For difficulty calculation we will only look at the highest strain value in each time interval of size STRAIN_STEP.
        // This is to eliminate higher influence of stream over aim by simply having more HitObjects with high strain.
        // The higher this value, the less strains there will be, indirectly giving long beatmaps an advantage.
        private const double STRAIN_STEP = 400;

        // The weighting of each strain value decays to 0.9 * it's previous value
        private const double DECAY_WEIGHT = 0.9;

        private double CalculateDifficulty(List<TpHitObject> tpHitObjects, DifficultyType Type)
        {
            // Find the highest strain value within each strain step
            List<double> HighestStrains = new List<double>();
            double IntervalEndTime = STRAIN_STEP;
            double MaximumStrain = 0; // We need to keep track of the maximum strain in the current interval

            TpHitObject PreviousHitObject = null;
            foreach (TpHitObject hitObject in tpHitObjects)
            {
                // While we are beyond the current interval push the currently available maximum to our strain list
                while (hitObject.BaseHitObject.StartTime > IntervalEndTime)
                {
                    HighestStrains.Add(MaximumStrain);

                    // The maximum strain of the next interval is not zero by default! We need to take the last hitObject we encountered, take its strain and apply the decay
                    // until the beginning of the next interval.
                    if (PreviousHitObject == null)
                    {
                        MaximumStrain = 0;
                    }
                    else
                    {
                        double Decay = Math.Pow(TpHitObject.DECAY_BASE[(int)Type], (double)(IntervalEndTime - PreviousHitObject.BaseHitObject.StartTime) / 1000);
                        MaximumStrain = PreviousHitObject.Strains[(int)Type] * Decay;
                    }

                    // Go to the next time interval
                    IntervalEndTime += STRAIN_STEP;
                }

                // Obtain maximum strain
                if (hitObject.Strains[(int)Type] > MaximumStrain)
                {
                    MaximumStrain = hitObject.Strains[(int)Type];
                }

                PreviousHitObject = hitObject;
            }

            // Build the weighted sum over the highest strains for each interval
            double Difficulty = 0;
            double Weight = 1;
            
            // Sort from highest to lowest strain.
            HighestStrains.Sort((a, b) => b.CompareTo(a));

            foreach (double Strain in HighestStrains)
            {
                Difficulty += Weight * Strain;
                Weight *= DECAY_WEIGHT;
            }

            return Difficulty;
        }
    }
}
