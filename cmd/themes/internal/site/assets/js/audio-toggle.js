// Audio toggle for the Uganda theme background music.
// Toggles play/pause on the anthem audio element and updates the button icon.

(function () {
	'use strict';

	var audio = document.getElementById('anthem');
	var toggle = document.getElementById('audio-toggle');
	var icon = document.querySelector('.audio-toggle__icon');

	if (!audio || !toggle || !icon) {
		return;
	}

	var isPlaying = false;

	function updateIcon() {
		icon.textContent = isPlaying ? '🔊' : '🔇';
		toggle.setAttribute('aria-label', isPlaying ? 'Pause background music' : 'Play background music');
	}

	// Start muted by default — browsers block autoplay with sound.
	// The user must click to unmute and play.
	audio.volume = 0.3;
	updateIcon();

	toggle.addEventListener('click', function () {
		if (isPlaying) {
			audio.pause();
			isPlaying = false;
		} else {
			audio.play().then(function () {
				isPlaying = true;
			}).catch(function () {
				// Autoplay blocked — show muted icon to indicate user action needed
				isPlaying = false;
			});
		}
		updateIcon();
	});

	// If audio ends (shouldn't happen with loop, but just in case)
	audio.addEventListener('ended', function () {
		isPlaying = false;
		updateIcon();
	});
})();
