-- Seed the 5 scoring categories (balanced weights).

INSERT INTO categories (slug, name, description, default_weight, display_order) VALUES
    ('economy', 'Economy', 'Employment, income, and economic growth indicators', 0.20, 1),
    ('education', 'Education', 'Educational attainment, outcomes, and investment', 0.20, 2),
    ('health', 'Health', 'Health outcomes, access, and risk factors', 0.20, 3),
    ('safety', 'Safety', 'Crime rates, traffic safety, and emergency response', 0.20, 4),
    ('affordability', 'Affordability', 'Housing costs, taxes, and cost of living', 0.20, 5);
