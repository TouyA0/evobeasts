package main

const krakenStage2Art = `       .--""--.__
     .'          ` + "`" + `--.__
    ( o                 ` + "`" + `--..__
     ` + "`" + `.                        ` + "`" + `>
       ` + "`" + `--.__          __.--'
             ` + "`" + `--.____.-'
           (          )
            (        )
             (      )
              (    )
               ` + "`" + `--'`
 
var krakenStage2Colors = [][3]int{
	{150, 180, 195},
	{150, 180, 195},
	{160, 188, 200},
	{165, 190, 200},
	{170, 195, 200},
	{175, 198, 200},
	{235, 150, 60},
	{235, 150, 60},
	{230, 140, 55},
	{225, 130, 50},
	{215, 120, 45},
}

func init() {
	speciesArts[EspeceKraken+"-2"] = CreatureArt{Art: krakenStage2Art, Colors: krakenStage2Colors}
}