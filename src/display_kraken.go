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

const krakenStage3Art = `           ,^.
       .--'  ` + "`" + `--.__
   .-''             ` + "`" + `--.__    \
 ,'                        ` + "`" + `-. >
( o                          .-'
 ` + "`" + `-._                  __.--'  /
     ` + "`" + `--..______..---''
      \  /
       \/`
 
var krakenStage3Colors = [][3]int{
	{80, 150, 160},
	{80, 150, 160},
	{90, 160, 165},
	{100, 170, 170},
	{140, 195, 185},
	{130, 185, 180},
	{110, 175, 175},
	{90, 150, 160},
	{80, 140, 155},
}

const krakenStage4Art = `          ,--^--. ,--^--,
         ((  @@  ))
    __,-'` + "`" + `==vv==` + "`" + `'-,__
  ,'   \\    ||    //   ` + "`" + `,
 /    __\\   ||   //__    \
|   ,'   ` + "`" + `. _||_ .'   ` + "`" + `.   |
 \ /       \/  \/       \ /
  X        /  \        X
 / \      (    )      / \
(   )      ` + "`" + `--'      (   )
~~~~~~~~~~~~~~~~~~~~~~~~~~~~`
 
var krakenStage4Colors = [][3]int{
	{60, 130, 120},
	{130, 240, 210},
	{50, 110, 100},
	{45, 100, 90},
	{40, 90, 82},
	{38, 85, 78},
	{35, 78, 72},
	{32, 72, 68},
	{28, 65, 62},
	{25, 58, 56},
	{225, 235, 235},
}

func init() {
	speciesArts[EspeceKraken+"-2"] = CreatureArt{Art: krakenStage2Art, Colors: krakenStage2Colors}
	speciesArts[EspeceKraken+"-3"] = CreatureArt{Art: krakenStage3Art, Colors: krakenStage3Colors}
	speciesArts[EspeceKraken+"-4"] = CreatureArt{Art: krakenStage4Art, Colors: krakenStage4Colors}
}
