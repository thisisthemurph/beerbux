import { useBackNavigationStore } from "@/stores/back-navigation-store.ts";
import { useEffect } from "react";
import { useSearchParams } from "react-router";

export const useBackNavigation = (link: string) => {
	const [searchParams] = useSearchParams();
	const setBackLink = useBackNavigationStore((state) => state.setBackLink);
	const clearBackLink = useBackNavigationStore((state) => state.clearBackLink);

	let backLink = link;

	const backLinkOverride = searchParams.get("bl");
	if (backLinkOverride) {
		backLink = backLinkOverride;
	}

	useEffect(() => {
		setBackLink(backLink);
		return () => clearBackLink();
	}, [backLink, setBackLink, clearBackLink]);
};

/**
 * Adds a back link override to the given useBackNavigation hook URL.
 * Returns the original urlPath if the backLinkOverride is not provided.
 * @param urlPath The URL path to add the back link override to.
 * @param backLinkOverride The back link override to add to the URL path.
 */
export function withBackLinkOverride(urlPath: string, backLinkOverride: string | undefined | null): string {
	if (!backLinkOverride) return urlPath;
	const urlEncodedPath = encodeURIComponent(backLinkOverride);
	return `${urlPath}?bl=${urlEncodedPath}`;
}
