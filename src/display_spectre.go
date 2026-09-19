package main

const spectreStage2Art = `      _,--._
    ,'      ` + "`" + `.
   ( (o)      \--.__
    ` + "`" + `.        ,-'   ` + "`" + `--.__
      ` + "`" + `--.__,-'           ` + "`" + `--.
      ,-'                 ,-. \
     (                   (   ) )
      ` + "`" + `-.                 ` + "`" + `-','
         ` + "`" + `-.__          _.-'
        (____) ` + "`" + `--....-'
         '  '`
 
var spectreStage2Colors = [][3]int{
	{95, 175, 105},
	{110, 190, 115},
	{125, 200, 120},
	{130, 205, 125},
	{120, 195, 120},
	{115, 190, 115},
	{110, 185, 115},
	{105, 180, 110},
	{95, 170, 105},
	{85, 155, 95},
	{75, 140, 85},
}

func init() {
	speciesArts[EspeceSpectre+"-2"] = CreatureArt{Art: spectreStage2Art, Colors: spectreStage2Colors}
}
