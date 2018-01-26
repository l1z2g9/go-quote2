class AudioTrack extends React.Component {
  constructor(props) {
      super(props);
      this.state = {url: "", title: ""};
      this.handleChange = this.handleChange.bind(this);
      this.playEnd = this.playEnd.bind(this);
   }

  handleChange(e) {
    this.url = "https://drive.google.com/uc?id=" + e.target.value + "&authuser=0&export=download";
    this.title = e.target.dataset.title;

    this.setState({url: this.url, title: this.title});
  }

  playEnd(e){
      let track = Math.floor((Math.random() * this.props.list.length) + 1);
      this.url = "https://drive.google.com/uc?id=" + this.props.list[track][0] + "&authuser=0&export=download";
      this.title = this.props.list[track][1];
      this.setState({url: this.url, title: this.title});
  }

  render() {
    return (
      <div>
         <table className="pure-table">
          <thead>
              <tr className="bg-info">
                  <td>Soundtrack</td>
                  <td>Play</td>
              </tr>
          </thead>
          <tbody>
            {
              this.props.list.map( (list, i) => {
                  return <tr className={i%2==0 ? 'pure-table-odd' : ''} key={list[0]}>
                        <td>
                {list.length > 2 &&
                          
                            <a href="https://drive.google.com/uc?id={list[2]}&authuser=0&export=download">{list[1]}</a>
                }
                {list.length ==2 &&
                            list[1]
                         
                }</td>
                    <td><button type="button" className="pure-button pure-button-primary" data-title={list[1]} onClick={this.handleChange} value={list[0]}>Play</button></td>

                </tr>
             })
            }
          </tbody>
          <tfoot>
              <tr>
                  <td colSpan="5" className="bg-info">
                      <b>{this.state.title}</b>
                  </td>
              </tr>
          </tfoot>
        </table>
        <div>
          <audio controls src={this.state.url} autoPlay onEnded={this.playEnd}></audio>
        </div>
      </div>
    )
  }
}