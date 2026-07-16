export type AppStatus = {
	status: string;
	version: string;
	databaseReady: boolean;
	activeImports: number;
	startedAt?: string;
};
export type State = { id: number; code: string; name: string; region?: string; division?: string };
export type Category = {
	id: number;
	slug: string;
	name: string;
	description?: string;
	defaultWeight: number;
	displayOrder: number;
};
export type Metric = {
	id: number;
	categoryId: number;
	slug: string;
	name: string;
	description?: string;
	unit?: string;
	higherIsBetter: boolean;
	normalizationMethod: string;
	defaultWeight: number;
	active?: boolean;
};
export type MetricValueQuality = {
	reportingCoverage?: number;
	participatingAgencies?: number;
	populationCovered?: number;
	dataRevision?: string;
	scoringEligible: boolean;
	exclusionReason?: string;
};
export type MetricValue = {
	id: number;
	stateId: number;
	metricId: number;
	year: number;
	value: number;
	sourceRecordId?: string;
	importId?: number;
	quality?: MetricValueQuality;
};
export type Profile = {
	id: number;
	name: string;
	description?: string;
	isDefault: boolean;
	isSystem: boolean;
};
export type CategoryWeight = { profileId: number; categoryId: number; weight: number };
export type CategoryScore = { categoryId: number; score: number; completeness: number };
export type StateScore = {
	stateId: number;
	overallScore: number;
	completeness: number;
	calculatedAt?: string;
	calculationVersion?: string;
	categories: CategoryScore[];
};
export type Scoreboard = {
	profileId: number;
	year: number;
	asOfYear: number;
	calculationVersion: string;
	scores: StateScore[];
};
export type DataSource = {
	id: number;
	name: string;
	publisher?: string;
	sourceUrl?: string;
	license?: string;
	format: string;
	description?: string;
	createdAt?: string;
	updatedAt?: string;
};
export type DataImport = {
	id: number;
	sourceId?: number;
	status: string;
	startedAt?: string;
	completedAt?: string;
	recordsRead: number;
	recordsInserted: number;
	recordsRejected: number;
	checksum?: string;
	errorSummary?: string;
};
export type ImportIssue = {
	id: number;
	importId: number;
	rowNumber?: number;
	fieldName?: string;
	rawValue?: string;
	errorMessage: string;
};
export type PublicSourceAdapter = {
	id: string;
	name: string;
	publisher: string;
	metricSlugs: string[];
	defaultYear: number;
	available: boolean;
	unavailableReason?: string;
};
export class ApiError extends Error {
	status: number;
	constructor(status: number, message: string) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
	}
}
