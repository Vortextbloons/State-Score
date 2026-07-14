import {
	getCategories,
	getMetrics,
	getScores,
	getStates,
	getValues,
	type Category,
	type Metric,
	type MetricValue,
	type Scoreboard,
	type State
} from '$lib/api';

export type Row = {
	state: State;
	overall: number | null;
	categories: Record<number, number | null>;
	completeness: number;
	values: MetricValue[];
};

export type ScoreData = {
	states: State[];
	categories: Category[];
	metrics: Metric[];
	rows: Row[];
	years: number[];
	asOfYear: number | null;
	calculationVersion: string | null;
	relative: true;
};

/** Load rankings from backend snapshots. Optional weights only re-average category scores. */
export async function loadScores(
	weights?: Record<number, number>,
	year?: number
): Promise<ScoreData> {
	const [states, allCategories, metrics] = await Promise.all([
		getStates(),
		getCategories(),
		getMetrics()
	]);
	const categories = allCategories.filter((c) => metrics.some((m) => m.categoryId === c.id));
	const everyValue = await getValues();
	const years = [...new Set(everyValue.map((v) => v.year))].sort((a, b) => b - a);

	let board: Scoreboard | null = null;
	try {
		board = await getScores(0, year ?? 0);
	} catch {
		board = null;
	}

	const asOfYear = board?.asOfYear ?? year ?? years[0] ?? null;
	const metricYears = new Map<number, number>();
	for (const m of metrics) {
		const available = [
			...new Set(
				everyValue
					.filter((v) => v.metricId === m.id && (asOfYear == null || v.year <= asOfYear))
					.map((v) => v.year)
			)
		].sort((a, b) => b - a);
		if (available[0]) metricYears.set(m.id, available[0]);
	}

	const byState = new Map(board?.scores.map((s) => [s.stateId, s]) ?? []);
	const rows: Row[] = states.map((state) => {
		const snap = byState.get(state.id);
		const cs: Record<number, number | null> = {};
		for (const c of categories) {
			const found = snap?.categories.find((x) => x.categoryId === c.id);
			cs[c.id] = found ? found.score : null;
		}
		let overall = snap?.overallScore ?? null;
		if (weights && overall != null) {
			const available = categories.filter((c) => cs[c.id] != null);
			const totalW = available.reduce((n, c) => n + (weights[c.id] ?? c.defaultWeight), 0);
			overall = totalW
				? available.reduce((n, c) => n + (cs[c.id] ?? 0) * (weights[c.id] ?? c.defaultWeight), 0) /
					totalW
				: null;
		}
		return {
			state,
			overall,
			categories: cs,
			completeness: snap?.completeness ?? 0,
			values: everyValue.filter(
				(v) => v.stateId === state.id && v.year === metricYears.get(v.metricId)
			)
		};
	});
	rows.sort((a, b) => (b.overall ?? -1) - (a.overall ?? -1));

	return {
		states,
		categories,
		metrics,
		rows,
		years,
		asOfYear,
		calculationVersion: board?.calculationVersion ?? null,
		relative: true
	};
}

export function fmt(v: number | null, digits = 1) {
	return v == null ? '—' : v.toFixed(digits);
}
