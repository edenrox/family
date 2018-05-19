convert -size 32x32 xc:none background.png

composite -geometry +1+6 state/US/AL.png background.png temp.png
composite overlay.png temp.png AL.png