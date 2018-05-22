convert -size 32x32 xc:none background.png

convert state/IE/flat/C.png -resize 30x20! flag-a.png
composite -geometry +1+6 flag-a.png background.png flag-b.png
composite overlay.png flag-b.png C.png
    
convert state/IE/flat/L.png -resize 30x20! flag-a.png
composite -geometry +1+6 flag-a.png background.png flag-b.png
composite overlay.png flag-b.png L.png
    
convert state/IE/flat/M.png -resize 30x20! flag-a.png
composite -geometry +1+6 flag-a.png background.png flag-b.png
composite overlay.png flag-b.png M.png
    
convert state/IE/flat/U.png -resize 30x20! flag-a.png
composite -geometry +1+6 flag-a.png background.png flag-b.png
composite overlay.png flag-b.png U.png
