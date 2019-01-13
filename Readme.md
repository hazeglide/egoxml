# EgoXML

This project is a rather hacky parser for x4:foundation savegames. the parser extracts the playerlog and saves it inside a sqlite database.the frontend uses those logs to generate income graphs per ship. 

### Config

``` toml
savegame = "C:/Users/<user>/Documents/Egosoft/X4/<id>/save/save_008.xml"
exclude = ["CJU-133","SBH-279","YSU-034","GLU-716","XYN-441","KMX-925"]

[combine]
SVF-965 = ["REF-314","IIB-734"]
SZS-033 = ["HIL-927","DXJ-278","LTN-250"]

[parser]
shipindex = 0
```
| config    	| description                                                                                                               	|
|-----------	|---------------------------------------------------------------------------------------------------------------------------	|
| savegame  	| path to the uncompressed savegame                                                                                         	|
| exclude   	| list of ship id's that will be ignored                                                                                    	|
| combine   	| mapping to combine ship id's, eg. station , build storage and station trader                                                  |
| shipindex 	| tells the parser on what position the current ship id can be found. 0 = first ship id in text, 1 = second, 2 = third, etc 	|
