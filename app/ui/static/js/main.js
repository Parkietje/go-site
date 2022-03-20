var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

document.getElementById("menu-button").addEventListener("click", toggleMenu);

function toggleMenu() {
	if(document.getElementById("sidebar").style.display == "none") {
		document.getElementById("sidebar").style.display = "block";
	} else {
		document.getElementById("sidebar").style.display = "none";
	}
} 