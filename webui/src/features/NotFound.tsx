import { Button } from "@/components/ui/button.tsx";
import { ArrowLeft, BeerOff } from "lucide-react";
import { useMemo } from "react";
import { Link } from "react-router";

const quotes = [
	"Beer not found. Please try again.",
	"Looks like the beer keg is empty here.",
	"If this page were a beer, it would be the bottom of the glass.",
	"All foam, no beer. Nothing to see here.",
	"You’ve reached the end of the keg.",
	"The good news: there's still beer somewhere. The bad news: not here.",
];

export default function NotFoundPage() {
	const quote = useMemo(() => quotes[Math.floor(Math.random() * quotes.length)], []);

	return (
		<div className="flex flex-col items-center justify-center">
			<BeerOff className="size-36 text-primary my-8" />
			<h1 className="text-4xl font-bold mb-8">Page not found</h1>
			<p className="mt-4 mb-8 text-2xl text-center tracking-wide">{quote}</p>
			<Button variant="secondary" asChild>
				<Link to="/">
					<ArrowLeft />
					<span>
						Stagger back <span className="font-semibold">Home</span>
					</span>
				</Link>
			</Button>
		</div>
	);
}
