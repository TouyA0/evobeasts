package main

const eggArt = `        .--""--.
      ,'        ` + "`" + `,
     /    .-.     \
    |    (   )     |
    |     ` + "`" + `-'      |
   |                 |
   |                 |
    \               /
     \             /
      ` + "`" + `.         ,'
        ` + "`" + `.     ,'
          ` + "`" + `---'`

var eggColors = [][3]int{
	{245, 235, 215},
	{248, 240, 222},
	{250, 245, 230},
	{252, 248, 236},
	{248, 240, 222},
	{243, 232, 210},
	{240, 228, 205},
	{236, 222, 198},
	{232, 216, 190},
	{228, 210, 182},
	{222, 202, 172},
	{215, 192, 160},
}

func init() {
	egg := CreatureArt{Art: eggArt, Colors: eggColors}
	speciesArts[EspecePhenix+"-1"] = egg
	speciesArts[EspeceKraken+"-1"] = egg
	speciesArts[EspeceSpectre+"-1"] = egg
}