import {
	ApiError,
	type AppStatus,
	type Category,
	type CategoryWeight,
	type DataImport,
	type DataSource,
	type ImportIssue,
	type Metric,
	type MetricValue,
	type Profile,
	type PublicSourceAdapter,
	type Scoreboard,
	type State
} from './types';

const BASE = '/api/v1';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const multipart = init?.body instanceof FormData;
	const options: RequestInit = {
		...init,
		headers: {
			Accept: 'application/json',
			...(!multipart ? { 'Content-Type': 'application/json' } : {}),
			...(init?.headers ?? {})
		}
	};
	let response: Response | null = null;
	for (let attempt = 0; attempt < 3; attempt++) {
		response = await fetch(BASE + path, options);
		if (![502, 503].includes(response.status) || (init?.method && init.method !== 'GET')) break;
		await new Promise((resolve) => setTimeout(resolve, 250 * (attempt + 1)));
	}
	if (!response!.ok) {
		let message = response!.statusText;
		try {
			const body = await response!.json();
			message = body.error?.message ?? message;
		} catch {
			/* non-JSON error */
		}
		throw new ApiError(response!.status, message);
	}
	return response!.json();
}

async function data<T>(path: string, init?: RequestInit): Promise<T> {
	return (await request<{ data: T }>(path, init)).data;
}

export const getStatus = () => request<AppStatus>('/status');
export const getStates = (region = '') =>
	data<State[]>('/states' + (region ? `?region=${encodeURIComponent(region)}` : ''));
export const getState = (code: string) => data<State>('/states/' + encodeURIComponent(code));
export const getCategories = () => data<Category[]>('/categories');
export const getMetrics = () => data<Metric[]>('/metrics');
export const getValues = (stateId?: number, year?: number) => {
	const q = new URLSearchParams();
	if (stateId) q.set('state_id', String(stateId));
	if (year) q.set('year', String(year));
	return data<MetricValue[]>('/values' + (q.size ? `?${q}` : ''));
};
export const getProfiles = () => data<Profile[]>('/profiles');
export const getDefaultProfile = () =>
	data<{ profile: Profile; categoryWeights: CategoryWeight[] }>('/profiles/default');
export const getSources = () => data<DataSource[]>('/sources');
export const saveSource = (source: Partial<DataSource>) =>
	data<DataSource>(source.id ? `/sources/${source.id}` : '/sources', {
		method: source.id ? 'PUT' : 'POST',
		body: JSON.stringify(source)
	});
export const getImports = () => data<DataImport[]>('/imports');
export const getImport = (id: number) =>
	data<{ import: DataImport; errors: ImportIssue[] }>(`/imports/${id}`);
export const uploadCSV = (sourceId: number, file: File) => {
	const body = new FormData();
	body.set('source_id', String(sourceId));
	body.set('file', file);
	return data<DataImport>('/imports', { method: 'POST', body });
};
export const getPublicSources = () => data<PublicSourceAdapter[]>('/public-sources');
export const refreshPublicSources = (adapterIds: string[], year?: number) =>
	data<{ imports: Record<string, number> }>('/public-sources/refresh', {
		method: 'POST',
		body: JSON.stringify({ adapterIds, year: year ?? 0 })
	});
export const getScores = (profileId = 0, year = 0) => {
	const q = new URLSearchParams();
	if (profileId) q.set('profile_id', String(profileId));
	if (year) q.set('year', String(year));
	return data<Scoreboard>('/scores' + (q.size ? `?${q}` : ''));
};
export const recalculate = (profileId = 0, year = 0) =>
	data<{ profileId: number; year: number; statesCalculated: number }>('/scores/recalculate', {
		method: 'POST',
		body: JSON.stringify({ profileId, year })
	});

export const api = {
	getStatus,
	getStates,
	getState,
	getCategories,
	getMetrics,
	getValues,
	getProfiles,
	getDefaultProfile,
	getSources,
	saveSource,
	getImports,
	getImport,
	uploadCSV,
	getPublicSources,
	refreshPublicSources,
	getScores,
	recalculate
};
