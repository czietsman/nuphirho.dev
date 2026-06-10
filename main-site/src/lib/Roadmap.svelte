<script lang="ts">
	interface Platform {
		label: string;
		url: string | null;
		pending?: boolean;
	}

	interface DayEntry {
		category: string;
		status: 'published' | 'pending-crosspost' | 'planned';
		title?: string;
		label?: string;
		'pending-title'?: string;
		platforms?: Platform[];
	}

	interface Week {
		label: string;
		days: Record<string, DayEntry>;
	}

	interface Month {
		month: string;
		label: string;
		future: boolean;
		weeks: Week[];
	}

	interface CalendarData {
		months: Month[];
	}

	let { data }: { data: CalendarData } = $props();

	const DAY_ORDER = ['mon', 'tue', 'wed', 'thu', 'fri', 'sat', 'sun'];
	const DAY_LABELS: Record<string, string> = { mon: 'Mon', tue: 'Tue', wed: 'Wed', thu: 'Thu', fri: 'Fri', sat: 'Sat', sun: 'Sun' };
	const DAY_CLASS: Record<string, string> = { mon: 'dark', tue: 'dark', wed: 'wed', thu: 'dark', fri: 'fri', sat: 'weekend', sun: 'weekend' };

	let collapsed = $state<Record<string, boolean>>(
		Object.fromEntries(data.months.map((m) => [m.month, m.future]))
	);

	function toggleMonth(month: string) {
		collapsed[month] = !collapsed[month];
	}
</script>

{#each data.months as month (month.month)}
	<div class="month-section">
		<div class="month-header" role="button" tabindex="0"
			onclick={() => toggleMonth(month.month)}
			onkeydown={(e) => e.key === 'Enter' && toggleMonth(month.month)}
		>
			<span>{month.label}</span>
			<span class="toggle">{collapsed[month.month] ? '▸ Show' : '▾ Hide'}</span>
		</div>

		{#if !collapsed[month.month]}
			<div class="month-content">
				{#each month.weeks as week}
					<div class="week">
						<div class="week-label">
							<div class="wk">Week</div>
							<div class="dates">{week.label}</div>
						</div>
						{#each DAY_ORDER as day}
							{@const entry = week.days?.[day]}
							<div class="day {DAY_CLASS[day]}">
								<div class="day-name">{DAY_LABELS[day]}</div>
								{#if entry}
									{@const catClass = entry.status === 'published' ? `published ${entry.category}`
										: entry.status === 'pending-crosspost' ? 'pending-crosspost'
										: `planned ${entry.category}`}
									<div class="post {catClass}">
										{#if entry.status === 'published'}<span class="check">✓</span>{/if}
										<span class="cat">{entry.label ?? entry.category}</span>
										{#if entry.title}
											<span class="title">{entry.title}</span>
										{:else}
											<span class="pending-title">Coming soon</span>
										{/if}
										{#if entry.platforms?.length}
											<div class="platforms">
												{#each entry.platforms as p}
													{#if p.pending || !p.url}
														<span class="platform-icon pending">{p.label}</span>
													{:else}
														<a class="platform-icon" href={p.url} target="_blank" rel="noopener">{p.label}</a>
													{/if}
												{/each}
											</div>
										{/if}
									</div>
								{/if}
							</div>
						{/each}
					</div>
				{/each}
			</div>
		{/if}
	</div>
{/each}
