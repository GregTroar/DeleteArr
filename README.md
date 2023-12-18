## DeleteArr

When hosting Radarr and Sonarr at home while using a remote seedbox for downloading creates an issue when you download your media back home from the seedbox.
Indeed after importing the files to your final media folder the files remains in your downloaded folder as well.
This small program removes the downloaded files after import.

## Install notes

go build DeleteArr.go

rename the config-samle.yml to config.yml and update as needed.

In order for the script to work properly you need to list down the folders where you download your stuff.
For example for me:

/mnt/Multimedia/Download/PostProcess/Movies
/mnt/Multimedia/Download/PostProcess/4K-Movies
/mnt/Multimedia/Download/PostProcess/Series
/mnt/Multimedia/Download/PostProcess/4K-Series
/mnt/Multimedia/Download/PostProcess/Kids

Add only the last folder in your config file as per the config-sample.yml

## Radarr / Sonarr config
Just go to settings / Connect / Custom Script and select the program and save
