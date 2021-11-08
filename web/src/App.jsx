import Header from "./components/Header";
import Footer from "./components/Footer";
import Endpoints from "./components/Endpoints";

function App() {
	return (
		<>
			<Header />
			<main className={"main"}>
				<Endpoints />
			</main>
			<Footer />
		</>
	);
}

export default App;
