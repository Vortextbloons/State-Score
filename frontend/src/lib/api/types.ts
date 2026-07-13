/** Shared API types for the local Go server. */

export type AppStatus = {
	status: string;
	version: string;
	databaseReady: boolean;
	activeImports: number;
	startedAt?: string;
};

export class ApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
	}
}
