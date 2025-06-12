-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW user_credit_score AS
WITH ledger_summary AS (
    SELECT
        user_id,
        SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END) AS beers_given,
        SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END) AS beers_received,
        -- Applies a time decay based on how many weeks old the transaction is.
        -- 6 days represents 0 weeks, 13 days represents 1 week etc.
        SUM(CASE WHEN amount > 0 THEN amount / (1 + FLOOR(EXTRACT(DAY FROM now() - created_at) / 7)) ELSE 0 END) AS recent_giving
    FROM ledger
    GROUP BY user_id
),
pairwise_giving AS (
    SELECT
        st.member_id AS giver_id,
        stl.member_id AS receiver_id,
        SUM(stl.amount) AS total_given
    FROM session_transactions st
    JOIN session_transaction_lines stl ON st.id = stl.transaction_id
    WHERE st.member_id <> stl.member_id
    GROUP BY st.member_id, stl.member_id
),
reciprocation AS (
    SELECT
        a.giver_id,
        a.receiver_id,
        a.total_given,
        COALESCE(b.total_given, 0) AS total_received_back,
        ROUND(COALESCE(b.total_given, 0) / NULLIF(a.total_given, 0), 2) AS reciprocation_ratio
    FROM pairwise_giving a
    LEFT JOIN pairwise_giving b
        ON a.giver_id = b.receiver_id AND a.receiver_id = b.giver_id
),
average_reciprocation AS (
    SELECT
        giver_id AS user_id,
        ROUND(AVG(reciprocation_ratio), 2) AS avg_reciprocation_ratio
    FROM reciprocation
    GROUP BY giver_id
),
combined AS (
    SELECT
        l.user_id,
        l.beers_given,
        l.beers_received,
        -- The ratio between beers given and received.
        ROUND(l.beers_given / (l.beers_received + 1), 2) AS balance_ratio,
        -- Number of beers given for every beer received.
        COALESCE(r.avg_reciprocation_ratio, 1.0) AS avg_reciprocation_ratio,
        ROUND(l.recent_giving, 2) AS recent_giving,
        -- Calculate the CreditScore
        -- Weighted score:
        --      60% balance_ratio
        --      30% reciprocation_ratio
        --      10% recency
        ROUND((l.beers_given / NULLIF(l.beers_received + 1, 0)) * 0.6 +
               COALESCE(r.avg_reciprocation_ratio, 1.0) * 0.3 +
               l.recent_giving * 0.1,
            2) * 100 AS credit_score
    FROM ledger_summary l
    LEFT JOIN average_reciprocation r ON l.user_id = r.user_id
),
scores AS (
    SELECT *,
       MIN(credit_score) OVER () AS min_credit_score,
       MAX(credit_score) OVER () AS max_credit_score
    FROM combined
)
SELECT
    user_id,
    beers_given,
    beers_received,
    balance_ratio,
    avg_reciprocation_ratio,
    recent_giving,
    ROUND((credit_score - min_credit_score) / NULLIF((max_credit_score - min_credit_score), 0) * 100, 2) AS credit_score,
    CASE
        WHEN (credit_score - min_credit_score) / NULLIF((max_credit_score - min_credit_score), 0) * 100 >= 80 THEN 'Round Champion'
        WHEN (credit_score - min_credit_score) / NULLIF((max_credit_score - min_credit_score), 0) * 100 >= 50 THEN 'Balanced Brewer'
        ELSE 'Round Dodger'
    END AS status_label
FROM scores;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS user_credit_score;
-- +goose StatementEnd
