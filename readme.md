# vcmanip
Split and Merge vCards

``vcmanip`` is unicode-friendly. So if you have a contact named "Íñigo Gámez",
his vCard will be "Íñigo Gámez.vcf"—instead of "igo_Gmez.vcf".

## Installation
````
go get github.com/pi241a/vcmanip
````

## Usage
````
Usage of vcmanip:
    -i="": vCard file or directory.
    -m=false: Merge a directory of vCards.
    -o="": Output directory.
    -s=false: Split a monolithic vCard.
````

Split multi-card vCard file into individual cards.
````
vcmanip -i big.vcf -s
````

Merge individual cards into one vCard file.
````
vcmanip -i contacts -m
````
