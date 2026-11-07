"""Underdog Woof by Conway Research (https://huggingface.co/ConwayResearch/Underdog-Woof).

The repository - the numpy runtime (the underdog_woof package) and the
40 MB int8 weights - is downloaded from Hugging Face into ctx.models_dir
on first load and imported from there. Woof makes 24 kHz mono audio.
"""

import os
import sys

from sonora import save_wav

REPO = "ConwayResearch/Underdog-Woof"
RATE = 24000


class Engine:
    def __init__(self, ctx):
        self.ctx = ctx
        self.voices = {}  # voice id -> Woof
        self.normalize = True

    def load(self, options):
        model_dir = os.path.join(self.ctx.models_dir, "Underdog-Woof")
        if not os.path.exists(os.path.join(model_dir, "woof_int8.npz")):
            from huggingface_hub import snapshot_download

            self.ctx.log("Downloading Underdog Woof (≈ 40 MB, first run only)…")
            snapshot_download(REPO, local_dir=model_dir, allow_patterns=["underdog_woof/*", "woof_int8.npz", "LICENSE", "NOTICE"])
        if model_dir not in sys.path:
            sys.path.insert(0, model_dir)
        self.normalize = bool(options.get("normalize", True))
        self.voices = {}
        self.woof("juno")
        return {}

    def woof(self, voice):
        """The model for voice, loaded the first time it's used."""
        if voice not in ("juno", "mira"):
            voice = "juno"
        if voice not in self.voices:
            from underdog_woof import Woof

            self.ctx.log(f"Loading Underdog Woof ({voice.title()})…")
            self.voices[voice] = Woof(voice=voice, normalize=self.normalize)
        return self.voices[voice]

    def synthesize(self, req):
        opts = req.get("options") or {}
        seed = str(opts.get("seed") or "").strip()
        audio = self.woof(req.get("voice") or "juno")(req["text"], seed=int(seed) if seed.isdigit() else None)
        return save_wav(req["output"], audio, RATE)
