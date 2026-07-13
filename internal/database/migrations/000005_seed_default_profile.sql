-- Seed the default Balanced scoring profile with category weights.

INSERT INTO scoring_profiles (name, description, is_default, is_system)
VALUES ('Balanced', 'Equal weight across all categories', 1, 1);

-- Add category weights for the Balanced profile (all 20%)
INSERT INTO profile_category_weights (profile_id, category_id, weight)
SELECT sp.id, c.id, 0.20
FROM scoring_profiles sp
CROSS JOIN categories c
WHERE sp.name = 'Balanced';
