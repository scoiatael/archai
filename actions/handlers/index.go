package handlers

import "github.com/scoiatael/archai/http"

const index = `
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">

  <title>Archai</title>
  <meta name="description" content="Simple page for interacting with Archai">
  <meta name="author" content="scoiatael">

  <!--[if lt IE 9]>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/html5shiv/3.7.3/html5shiv.js"></script>
  <![endif]-->
</head>

<body>
	<div id='app'>
	<h2> React is loading... </h2>
	</div>
	<script src="https://fb.me/react-15.0.2.js"></script>
	<script src="https://fb.me/react-dom-15.0.2.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/babel-core/5.8.23/browser.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.2.0/jquery.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/marked/0.3.5/marked.min.js"></script>
	<script type="text/babel">
	const Streams = ({streams, onClick}) => {
		let items = streams.map((stream) => {
		    let click = () => { onClick(stream) };
			return <li onClick={click} key={stream}>{stream}</li>
		});
		return <ul>{items}</ul>;
	}

	class Stream extends React.Component {
		constructor(props) {
			super(props);
			this.name = props.name;
			this.state = { "results": [] };
		}

		componentDidMount() {
			$.getJSON('/stream/' + this.name, (data) => {
				this.setState(data);
			});
		}

		render() {
			let items = this.state.results.map((result, index) =>
				<li key={index}>{JSON.stringify(result)}</li>
			);
			return <ul>{items}</ul>;
		}
	}

	class App extends React.Component {
		constructor(props) {
			super(props);
			this.state = { "streams": [], "screen": "streams" };
		}

		componentDidMount() {
			$.getJSON('/streams', (data) => {
				this.setState(data);
			});
		}

		setStream(name) {
			this.setState({ "screen": "stream", "stream": name});
		}

		render() {
			if(this.state.screen == "stream") {
				return <Stream name={this.state.stream} />;
			}
			return <Streams streams={this.state.streams} onClick={this.setStream.bind(this)} />;
		}
	}

	ReactDOM.render(
		<App />,
		document.getElementById('app')
	);
	</script>
</body>
</html>
`

func (gs Handler) Index(ctx http.GetContext) {
	ctx.SendHtml(index)
}
