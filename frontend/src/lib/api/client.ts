import { ApiError, type AppStatus } from './types';

const API_BASE = '/api/v1';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const response = await fetch(`${API_BASE}${path}`, {
		...init,
		headers: {
			Accept: 'application/json',
			...(init?.headers ?? {})
		}
	});

	if (!response.ok) {
		let message = response.statusText || 'Request failed';
		try {
			const body = (await response.json()) as { error?: string; message?: string };
			message = body.error ?? body.message ?? message;
		} catch {
			// keep statusText
		}
		throw new ApiError(response.status, message);
	}

	return (await response.json()) as T;
}

/** GET /api/v1/status — application readiness. */
export function getStatus(): Promise<AppStatus> {
	return request<AppStatus>('/status');
}

export const api = {
	getStatus
};
