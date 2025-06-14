export type CreditScoreStatus = "Round Dodger" | "Balanced Brewer" | "Round Champion";

const BalancedBrewer = 50;
const RoundChampion = 80;

export const getCreditScoreStatus = (score: number): CreditScoreStatus => {
	return score < BalancedBrewer
		? "Round Dodger"
		: score < RoundChampion
			? "Balanced Brewer"
			: "Round Champion";
};
