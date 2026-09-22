package main

const bannerMenu = "assets/banners/banniere_menu.png"
const bannerMerchant = "assets/banners/banniere_marchand.png"
const bannerForge = "assets/banners/banniere_forge.png"
const bannerWidth = 120

func displayBanner(path string) {
	displayImage(path, bannerWidth)
}