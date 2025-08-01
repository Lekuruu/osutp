﻿using System;

namespace osutp.TomPoints
{
    [Flags]
    public enum Mods
    {
        None                   = 0,
        NoFail                 = 1,
        Easy                   = 2,
        NoVideo                = 4,
        Hidden                 = 8,
        HardRock               = 16,
        SuddenDeath            = 32,
        DoubleTime             = 64,
        Relax                  = 128,
        HalfTime               = 256,
        Nightcore              = 512,
        Flashlight             = 1024,
        Autoplay               = 2048,
        SpunOut                = 4096,
        Relax2                 = 8192,
        Perfect                = 16384,
        Key4                   = 32768,
        Key5                   = 65536,
        Key6                   = 131072,
        Key7                   = 262144,
        Key8                   = 524288,
        FadeIn                 = 1048576,
        Random                 = 2097152,
        LastMod                = 4194304,
        KeyMod                 = Key4 | Key5 | Key6 | Key7 | Key8,
        FreeModAllowed         = NoFail | Easy | Hidden | HardRock | SuddenDeath | Flashlight | FadeIn | Relax | Relax2 | SpunOut | KeyMod
    }
}