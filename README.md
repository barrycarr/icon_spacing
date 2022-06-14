# Icon Spacing

A Windows 11 command line utility that will enable you to adjust the spacing of the icons on your desktop.

Usage:
    `icon_spacing -distance=<wide|medium|narrow> [-update]`

Where:

- `distance` sets the desired space between your icons. The options are:
  - `wide` - The default spacing used by Windows 11
  - `medium` - 'nuff said
  - `narrow` - The spacing used by previous versions of windows.
- `update` sets the new spacing without asking for confirmation.

Example:
    `icon_spacing -distance=narrow -update`

This example sets your icon spacing used by Windows before version 11. The `update` switch means that your registry will be updated with prompting you, first.

  
### Notes
- Only works on Windows 11 or greater. If a lower version of windows is detected the program will exit with a message.
- You will need to either log-out or reboot after you have used this utility before you see any effects

USE AT YOUR OWN RISK. BACKUP YOUR REGISTRY BEFORE ATTEMPTING TO USE THIS UTILITY. I ACCEPT NO LIABILITY IF YOU EXPERIENCE PROBLEMS.    
