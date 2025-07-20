import Navbar from "./components/Navigation";
import NotificationContainer from "./components/Notification";

function App(props: { children: Element }) {
	return (
		<div class="flex flex-col items-center min-h-screen bg-gray-100 w-screen">
            <Navbar/>

            {props.children}

            <NotificationContainer />
		</div>
	);
}

export default App;
