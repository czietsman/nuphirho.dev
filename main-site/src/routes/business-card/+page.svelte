<script lang="ts">
	import qrCodeUrl from './linkedin-qr.svg?url';
	import { BUSINESS_CARD, shareBusinessCard } from './card.js';

	let shareStatus = $state('');

	async function shareCard() {
		shareStatus = '';
		try {
			const result = await shareBusinessCard(navigator);
			shareStatus = result === 'shared' ? '' : 'Card link copied';
		} catch (error) {
			if (error instanceof DOMException && error.name === 'AbortError') return;
			shareStatus = 'Sharing is unavailable';
		}
	}
</script>

<svelte:head>
	<title>{BUSINESS_CARD.name} | Digital Business Card</title>
	<meta name="description" content="Contact and professional profile for Christo Zietsman" />
	<meta name="theme-color" content="#041619" />
	<meta property="og:title" content="Christo Zietsman | Digital Business Card" />
	<meta property="og:description" content="Contact Christo Zietsman at nuphirho.dev" />
	<meta property="og:url" content={BUSINESS_CARD.publicUrl} />
</svelte:head>

<div class="card-screen">
	<div class="water-light water-light-one"></div>
	<div class="water-light water-light-two"></div>

	<section class="card" aria-label="Digital business card for {BUSINESS_CARD.name}">
		<div class="paperless" aria-label="Paperless card" title="Paperless card">
			<svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
				<path d="M19.5 4.5C12.8 4.8 7.3 7.1 5.2 11.3c-1.5 3 .1 6.2 3.1 6.8 4.3.8 8.3-3.1 11.2-13.6Z" />
				<path d="M4.5 20c2.6-5.3 6.3-8.4 11.7-10.7" />
			</svg>
		</div>

		<div class="monogram" aria-hidden="true"><span>N</span></div>

		<div class="identity">
			<div class="details">
				<p class="eyebrow">Digital business card</p>
				<h1>{BUSINESS_CARD.name}</h1>
				<a class="contact-line" href="https://{BUSINESS_CARD.website}">{BUSINESS_CARD.website}</a>
				<a class="contact-line" href="tel:{BUSINESS_CARD.internationalPhone}">{BUSINESS_CARD.phone}</a>

				<div class="actions">
					<button type="button" onclick={shareCard}>
						<svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
							<circle cx="18" cy="5" r="2.5" />
							<circle cx="6" cy="12" r="2.5" />
							<circle cx="18" cy="19" r="2.5" />
							<path d="m8.2 10.8 7.6-4.5M8.2 13.2l7.6 4.5" />
						</svg>
						<span><strong>Share card</strong><small>Quick Share / AirDrop</small></span>
					</button>
					<a href="/business-card/contact.vcf" download>
						<svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
							<circle cx="12" cy="8" r="3" />
							<path d="M6.5 19c.7-3.2 2.5-4.8 5.5-4.8s4.8 1.6 5.5 4.8M19 9v6M16 12h6" />
						</svg>
						<span><strong>Add contact</strong><small>Android / Apple</small></span>
					</a>
				</div>
				{#if shareStatus}<p class="share-status" aria-live="polite">{shareStatus}</p>{/if}
			</div>

			<a class="qr-link" href={BUSINESS_CARD.linkedInUrl} target="_blank" rel="noopener noreferrer" aria-label="Open {BUSINESS_CARD.name} on LinkedIn">
				<span class="qr-frame"><img src={qrCodeUrl} alt="LinkedIn QR code" /></span>
				<span class="scan-label">Scan for LinkedIn</span>
			</a>
		</div>
	</section>
	<p class="signature">Nuphirho Research</p>
</div>

<style>
	.card-screen {
		position: relative;
		isolation: isolate;
		display: grid;
		min-height: 100dvh;
		width: 100%;
		overflow: hidden;
		place-items: center;
		padding: clamp(1rem, 4vw, 4rem);
		color: #eefcf8;
		background:
			radial-gradient(circle at 72% 18%, rgba(53, 169, 164, 0.16), transparent 30%),
			radial-gradient(circle at 18% 82%, rgba(21, 92, 101, 0.24), transparent 34%),
			linear-gradient(138deg, #02090b 0%, #041619 48%, #071f22 100%);
	}

	.card-screen::before {
		position: absolute;
		z-index: -2;
		inset: -35%;
		content: '';
		opacity: 0.24;
		background:
			repeating-radial-gradient(ellipse at 40% 48%, transparent 0 34px, rgba(105, 224, 211, 0.14) 36px 38px, transparent 42px 68px),
			repeating-radial-gradient(ellipse at 76% 28%, transparent 0 51px, rgba(117, 210, 201, 0.1) 53px 55px, transparent 59px 91px);
		filter: blur(1px);
		transform: rotate(-8deg) scale(1.1);
		animation: water-drift 18s ease-in-out infinite alternate;
	}

	.card-screen::after {
		position: absolute;
		z-index: -1;
		inset: 0;
		content: '';
		pointer-events: none;
		background: linear-gradient(108deg, transparent 12%, rgba(176, 255, 245, 0.035) 42%, transparent 66%);
		animation: light-sweep 12s ease-in-out infinite alternate;
	}

	.water-light {
		position: absolute;
		z-index: -1;
		border: 1px solid rgba(112, 220, 207, 0.1);
		border-radius: 50%;
		box-shadow: 0 0 80px rgba(64, 188, 178, 0.07), inset 0 0 60px rgba(64, 188, 178, 0.05);
	}

	.water-light-one { width: 54vw; height: 18vw; top: 8%; left: 26%; transform: rotate(-13deg); }
	.water-light-two { width: 46vw; height: 15vw; right: -8%; bottom: 10%; transform: rotate(17deg); }

	.card {
		position: relative;
		display: grid;
		grid-template-columns: minmax(170px, 0.78fr) minmax(0, 1.35fr);
		width: min(1120px, 100%);
		min-height: min(620px, calc(100dvh - clamp(2rem, 8vw, 8rem)));
		border: 1px solid rgba(169, 231, 222, 0.17);
		border-radius: clamp(1.5rem, 3vw, 3.25rem);
		background: linear-gradient(120deg, rgba(3, 17, 20, 0.88), rgba(8, 35, 38, 0.7));
		box-shadow: 0 35px 100px rgba(0, 0, 0, 0.46), inset 0 1px rgba(255, 255, 255, 0.05);
		backdrop-filter: blur(18px);
		overflow: hidden;
		animation: card-arrive 900ms cubic-bezier(0.16, 1, 0.3, 1) both;
	}

	.paperless {
		position: absolute;
		z-index: 3;
		top: clamp(1rem, 2.4vw, 1.8rem);
		right: clamp(1rem, 2.4vw, 1.8rem);
		display: grid;
		width: 34px;
		height: 34px;
		place-items: center;
		border: 1px solid rgba(91, 214, 128, 0.26);
		border-radius: 50%;
		color: #5bd680;
		background: rgba(16, 73, 44, 0.22);
		box-shadow: 0 0 24px rgba(71, 210, 113, 0.12);
	}

	.paperless svg { width: 19px; height: 19px; stroke: currentColor; stroke-width: 1.7; stroke-linecap: round; stroke-linejoin: round; }

	.monogram {
		position: relative;
		display: grid;
		place-items: center;
		min-width: 0;
		border-right: 1px solid rgba(169, 231, 222, 0.16);
		background: linear-gradient(155deg, rgba(54, 151, 145, 0.1), rgba(0, 0, 0, 0.16));
	}

	.monogram::after {
		position: absolute;
		width: min(72%, 230px);
		aspect-ratio: 1;
		border: 1px solid rgba(130, 226, 213, 0.14);
		border-radius: 50%;
		content: '';
		box-shadow: 0 0 80px rgba(65, 205, 190, 0.08);
	}

	.monogram span {
		position: relative;
		z-index: 1;
		font-family: 'Bodoni 72', Didot, 'Times New Roman', serif;
		font-size: clamp(9rem, 23vw, 18rem);
		font-weight: 400;
		line-height: 0.7;
		letter-spacing: -0.09em;
		color: #e9faf5;
		text-shadow: 0 0 38px rgba(129, 230, 216, 0.16);
		transform: translateX(-0.04em);
	}

	.identity {
		display: grid;
		grid-template-columns: minmax(0, 1fr) auto;
		align-items: center;
		gap: clamp(2rem, 5vw, 5.5rem);
		padding: clamp(2rem, 6vw, 6rem);
	}

	.details { display: flex; min-width: 0; flex-direction: column; align-items: flex-start; }
	.eyebrow { margin: 0 0 1.25rem; font: 600 0.7rem/1.2 'Avenir Next', Avenir, 'Century Gothic', sans-serif; letter-spacing: 0.24em; text-transform: uppercase; color: #79b9b1; }
	h1 { max-width: 12ch; margin: 0 0 1.4rem; font: 400 clamp(2.75rem, 6vw, 5.8rem)/0.92 'Bodoni 72', Didot, 'Times New Roman', serif; letter-spacing: -0.045em; color: #f2fbf8; }
	.contact-line { min-height: 40px; font: 400 clamp(1rem, 1.7vw, 1.35rem)/40px 'Avenir Next', Avenir, 'Century Gothic', sans-serif; letter-spacing: 0.035em; color: #b9d9d4; text-decoration: none; }
	.contact-line:hover { color: #fff; }

	.actions { display: flex; flex-wrap: wrap; gap: 0.65rem; margin-top: 1.35rem; }
	.actions button, .actions a { display: flex; min-height: 48px; align-items: center; gap: 0.65rem; padding: 0.55rem 0.85rem; border: 1px solid rgba(151, 218, 208, 0.2); border-radius: 0.75rem; color: #dff4ef; background: rgba(110, 194, 183, 0.07); font-family: 'Avenir Next', Avenir, 'Century Gothic', sans-serif; text-align: left; text-decoration: none; cursor: pointer; }
	.actions button:hover, .actions a:hover { border-color: rgba(151, 218, 208, 0.42); background: rgba(110, 194, 183, 0.13); }
	.actions svg { width: 20px; height: 20px; flex: 0 0 auto; stroke: #83c9bf; stroke-width: 1.6; stroke-linecap: round; stroke-linejoin: round; }
	.actions span { display: flex; flex-direction: column; }
	.actions strong { font-size: 0.76rem; font-weight: 600; }
	.actions small { margin-top: 0.05rem; font-size: 0.56rem; letter-spacing: 0.04em; color: #79aaa4; }
	.share-status { margin: 0.55rem 0 0; font: 500 0.68rem/1.2 'Avenir Next', Avenir, sans-serif; color: #8fd8ce; }

	.qr-link { display: flex; flex-direction: column; align-items: center; gap: 0.85rem; color: #8dc9c0; text-decoration: none; }
	.qr-frame { display: block; width: clamp(150px, 18vw, 210px); padding: clamp(0.65rem, 1.2vw, 0.9rem); border: 1px solid rgba(154, 231, 220, 0.32); border-radius: 1.2rem; background: #f3fbf8; box-shadow: 0 18px 55px rgba(0, 0, 0, 0.3), 0 0 40px rgba(64, 188, 178, 0.08); transition: transform 180ms ease; }
	.qr-link:hover .qr-frame { transform: translateY(-3px); }
	.qr-frame img { display: block; width: 100%; height: auto; }
	.scan-label { font: 600 0.67rem/1.2 'Avenir Next', Avenir, 'Century Gothic', sans-serif; letter-spacing: 0.2em; text-transform: uppercase; }
	.signature { position: absolute; right: clamp(1.25rem, 4vw, 4rem); bottom: clamp(1rem, 2.5vw, 2rem); margin: 0; font: 400 0.62rem/1.2 'Avenir Next', Avenir, sans-serif; letter-spacing: 0.19em; text-transform: uppercase; color: rgba(174, 218, 211, 0.5); }

	@keyframes card-arrive { from { opacity: 0; transform: translateY(18px) scale(0.985); } to { opacity: 1; transform: translateY(0) scale(1); } }
	@keyframes water-drift { from { transform: rotate(-8deg) scale(1.1) translate3d(-2%, -1%, 0); } to { transform: rotate(-4deg) scale(1.16) translate3d(3%, 2%, 0); } }
	@keyframes light-sweep { from { transform: translateX(-8%); opacity: 0.5; } to { transform: translateX(8%); opacity: 1; } }

	@media (max-width: 760px) {
		.card-screen { padding: 0; }
		.card { grid-template-columns: minmax(105px, 30vw) minmax(0, 1fr); min-height: 100dvh; border: 0; border-radius: 0; }
		.identity { grid-template-columns: 1fr; align-content: center; gap: clamp(1.7rem, 5vh, 3rem); padding: clamp(1.25rem, 5vw, 2.5rem); }
		.monogram span { font-size: clamp(6rem, 28vw, 10rem); }
		h1 { font-size: clamp(2.25rem, 10vw, 4rem); }
		.qr-frame { width: clamp(122px, 34vw, 164px); }
		.qr-link { align-items: flex-start; }
		.actions { margin-top: 1rem; }
		.signature { display: none; }
	}

	@media (max-height: 560px) and (orientation: landscape) {
		.card-screen { height: 100dvh; min-height: 0; padding: 0.65rem; overflow: hidden; }
		.card { height: calc(100dvh - 1.3rem); min-height: 0; grid-template-columns: minmax(140px, 0.65fr) minmax(0, 1.35fr); }
		.identity { gap: 1.35rem; padding: 0.85rem 2rem; }
		.eyebrow { margin-bottom: 0.45rem; }
		h1 { margin-bottom: 0.55rem; font-size: clamp(2.1rem, 5.5vw, 3rem); }
		.contact-line { min-height: 32px; font-size: 0.9rem; line-height: 32px; }
		.actions { gap: 0.45rem; margin-top: 0.7rem; }
		.actions button, .actions a { min-height: 42px; padding: 0.4rem 0.65rem; }
		.qr-frame { width: 102px; padding: 0.5rem; border-radius: 0.8rem; }
		.qr-link { gap: 0.45rem; }
		.paperless { top: 0.75rem; right: 0.75rem; width: 30px; height: 30px; }
		.signature { display: none; }
	}

	@media (prefers-reduced-motion: reduce) {
		.card, .card-screen::before, .card-screen::after { animation: none; }
	}
</style>
