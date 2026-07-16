import { describe, expect, it } from 'vitest';
import { formatPopulation } from './scores';

describe('formatPopulation', () => {
	it('uses compact M and k suffixes', () => {
		expect(formatPopulation(1_000_000)).toBe('1M');
		expect(formatPopulation(39_355_309)).toBe('39.4M');
		expect(formatPopulation(100_000)).toBe('100k');
		expect(formatPopulation(737_270)).toBe('737k');
	});
});
