package main

const bannerMenu = "assets/ui/banners/banniere_menu.png"
const bannerMerchant = "assets/ui/banners/banniere_marchand.png"
const bannerForge = "assets/ui/banners/banniere_forge.png"
const bannerInventaire = "assets/ui/banners/banniere_inventaire.png"
const bannerCombat = "assets/ui/banners/banniere_combat.png"
const bannerMemoires = "assets/ui/banners/banniere_memoires.png"
const bannerWidth = 120

func displayBanner(path string) {
	displayImage(path, bannerWidth)
}