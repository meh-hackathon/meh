import Navbar from "./components/Navigation";
import NotificationContainer from "./components/Notification";

function App() {

	return (
		<div class="flex flex-col items-center min-h-screen bg-gray-100 w-screen">
            <Navbar/>

            <h1 class="text-3xl font-bold underline">
                Hello world!
            </h1>

            <NotificationContainer />
		</div>
	);
}

export default App;
