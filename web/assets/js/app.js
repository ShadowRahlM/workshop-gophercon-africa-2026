console.log('🇺🇬 Uganda Wildlife Sightings');

// Background audio player
(function () {
	const audio = document.getElementById('bgAudio');
	const playBtn = document.getElementById('playBtn');
	const volumeSlider = document.getElementById('volumeSlider');
	const eqBars = document.getElementById('eqBars');

	if (!audio || !playBtn) return;

	// Set initial volume
	audio.volume = (volumeSlider.value || 30) / 100;
	eqBars.classList.add('is-paused');

	// Play/Pause toggle
	playBtn.addEventListener('click', function () {
		if (audio.paused) {
			audio.play().then(function () {
				playBtn.innerHTML = '<i class="ph ph-pause"></i>';
				eqBars.classList.remove('is-paused');
			}).catch(function (err) {
				console.warn('Audio play failed:', err);
			});
		} else {
			audio.pause();
			playBtn.innerHTML = '<i class="ph ph-play"></i>';
			eqBars.classList.add('is-paused');
		}
	});

	// Volume control
	if (volumeSlider) {
		volumeSlider.addEventListener('input', function () {
			audio.volume = this.value / 100;
		});
	}

	// Reset button when audio ends (shouldn't happen with loop, but safety)
	audio.addEventListener('ended', function () {
		playBtn.innerHTML = '<i class="ph ph-play"></i>';
		eqBars.classList.add('is-paused');
	});
})();

// 3D card tilt effect on mouse move
(function () {
	if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return;

	document.addEventListener('mousemove', function (e) {
		var cards = document.querySelectorAll('.card');
		cards.forEach(function (card) {
			var rect = card.getBoundingClientRect();
			var x = e.clientX - rect.left;
			var y = e.clientY - rect.top;

			// Only apply tilt if mouse is near the card
			if (x >= -50 && x <= rect.width + 50 && y >= -50 && y <= rect.height + 50) {
				var centerX = rect.width / 2;
				var centerY = rect.height / 2;
				var rotateX = ((y - centerY) / centerY) * -3;
				var rotateY = ((x - centerX) / centerX) * 3;
				card.style.transform = "translateY(-4px) rotateX(" + rotateX + "deg) rotateY(" + rotateY + "deg)";
			} else {
				card.style.transform = '';
			}
		});
	});
})();
