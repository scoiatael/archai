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
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.98.0/css/materialize.min.css" integrity="sha256-mDRlQYEnF3BuKJadRTD48MaEv4+tX8GVP9dEvjZRv3c=" crossorigin="anonymous" />

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
	<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.98.0/js/materialize.min.js" integrity="sha256-abuKx2bTKkpnebr/MelhYjv6tAZvfBQ2VKxpi2yJ57o=" crossorigin="anonymous"></script>
	<script src="https://unpkg.com/react-router-dom/umd/react-router-dom.min.js"></script>
	<script type="text/babel">
	const { HashRouter, Route, Link } = ReactRouterDOM

	const Streams = ({streams}) => {
		let items = streams.map((stream) => {
			let to = "/stream/" + stream
			return <Link to={to} className="collection-item" key={stream}>{stream}</Link>
		});
		return(
			<div className="container">
				<div className="row">
					<ul className="collection">{items}</ul>
				</div>
			</div>);
	}

	class Stream extends React.Component {
		constructor(props) {
			super(props);
			this.updateProps(props);
			this.state = { "results": [] };
		}

		componentWillReceiveProps(props) {
			this.updateProps(props);
			this.fetchData();
		}

		updateProps(props) {
			this.name = props.match.params.name;
			this.cursor = props.match.params.cursor || '';
		}

		fetchData() {
			$.getJSON('/stream/' + this.name + "?cursor=" + this.cursor, (data) => {
				this.setState(data);
			});
		}

		componentDidMount() {
			this.fetchData();
		}

		render() {
			let items = this.state.results.map((result, index) =>
				<li className="collection-item" key={index}>{JSON.stringify(result)}</li>
			);
			let more = '';
			let cursor = this.state.cursor && this.state.cursor.next;
			if (cursor) {
				more = <Link to={"/stream/" + this.name + "/" + cursor }>More</Link>;
			}
			return(
				<div className="container">
					<div className="row">
						<h3>{this.name}</h3>
						<ul className="collection">{items}</ul>
					</div>
					{ more }
				</div>
			);
		}
	}

	class App extends React.Component {
		constructor(props) {
			super(props);
			this.state = { "streams": [] };
		}

		componentDidMount() {
			$.getJSON('/streams', (data) => {
				this.setState(data);
			});
		}

		render() {
			return <Streams streams={this.state.streams} />;
		}
	}

	ReactDOM.render(
		<HashRouter>
			<div>
				<Route exact path="/" component={App}/>
				<Route exact path="/stream/:name" component={Stream}/>
				<Route path="/stream/:name/:cursor" component={Stream}/>
			</div>
		</HashRouter>,
		document.getElementById('app')
	);
	</script>
</body>
</html>
`

func (gs Handler) Index(ctx http.GetContext) {
	ctx.SendHtml(index)
}
