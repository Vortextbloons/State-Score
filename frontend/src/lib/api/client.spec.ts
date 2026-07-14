import { afterEach, describe, expect, it, vi } from 'vitest';
import { getValues } from './client';

describe('API client', () => {
	afterEach(() => vi.unstubAllGlobals());

	it('uses the bulk values endpoint when no state is supplied', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValue(
				new Response(JSON.stringify({ data: [] }), {
					status: 200,
					headers: { 'Content-Type': 'application/json' }
				})
			);
		vi.stubGlobal('fetch', fetchMock);
		await getValues();
		expect(fetchMock).toHaveBeenCalledWith('/api/v1/values', expect.any(Object));
	});
});
