# gtm - global theme manager

gtm is a small CLI tool that allows for configuration of various different packages via a single command or theme file.

The program works using inserts which are placed next to the config files of the respective tools in their .config dir and then included via a tool dependant include statement (for example source = ~/.config/hypr/gtm_hyprland.conf in hyprland.conf).

The gtm config is placed in .config/gtm/ as gtm.yaml and the themes are in a themes directory next to the config by default.
Samples for inserts, themes and a config are provided in the "config" dir in this repo.
The inserts basically just find and replace whatever is between the start and end marker in the comment at the top of the insert file. You may use the sample inserts as a baseline and copy your own config into it adhering to the general format given by the samples.



