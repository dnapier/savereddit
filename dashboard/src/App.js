import './App.css';
import React from 'react';
import {ArwesThemeProvider, Button, Figure, FrameCorners, FrameHexagon, StylesBaseline, Text} from "@arwes/core";
import {BleepsProvider, BleepsAudioSettings, BleepsPlayersSettings, BleepsSettings} from "@arwes/sounds";
import {Animator, AnimatorGeneralProvider} from "@arwes/animation";

const IMAGE_URL = 'https://playground.arwes.dev/assets/images/wallpaper.jpg';
const SOUND_OBJECT_URL = 'https://playground.arwes.dev/assets/sounds/object.mp3';
const SOUND_TYPE_URL = 'https://playground.arwes.dev/assets/sounds/type.mp3';

const audioSettings: BleepsAudioSettings = { common: { volume: 0.25, disabled: false }};
const playersSettings: BleepsPlayersSettings = { object: { src: [SOUND_OBJECT_URL] }, type: { src: [SOUND_TYPE_URL], loop: true}};
const bleepsSettings: BleepsSettings = { object: { player: 'object' }, type: { player: 'type' }};
const generalAnimator = { duration: { enter: 200, exit: 200} };
const url = 'http://127.0.0.1:8080/posts'

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      account: null,
      message: null,
      activate: true,
      text: null,
    }

    // Optional 1:
    // this.setAssimilate = this.setAssimilate.bind(this)
    // onClick={this.setAssimilate}
    // onClick={this.setAssimilate.bind(this, id)} // with args
    // Option 2:
    // onClick={() => this.setAssimilate()}
    // onClick={(e) => this.setAssimilate(id, e)} // with args

    this.setOnClick = this.setOnClick.bind(this)
  }

  setOnClick(m) {
    return () => this.setState({message: m})
  }

  componentDidMount() {}

  async setAssimilate() {
    // TODO - 2 - Retrieve all posts
    await fetch(url).then(response => {
      if (response.ok) {
        console.log(`response:`+response)
        return response.json()
      }
      throw response})
    .then(data => {this.setState({message: data.Attrs.Author, text: data.Attrs.SelfTextHTML})})
    .catch(error => console.error())
  }

  render() {
    let text = this.state.text;
    let activate = this.state.activate;
    let cornerDiv = <div style={{ width: 150, height: 25 }}>Blah</div>

    return (
      <div className="App" >
          <ArwesThemeProvider>
            <StylesBaseline styles={{ body: { fontFamily: '"Titillium Web", sans-serif' } }} />
            <BleepsProvider audioSettings={audioSettings} playersSettings={playersSettings} bleepsSettings={bleepsSettings}>
              <AnimatorGeneralProvider animator={generalAnimator}>
                <Animator animator={{ activate, manager: 'stagger' }}>
                  <nav>
                  {/*<nav style={{ display: 'flex', 'flex-direction': 'column' }}>*/}
                    <FrameCorners onClick={this.setOnClick("Strategery")} animator={{activate}} hover>{ cornerDiv }</FrameCorners>
                    <FrameCorners onClick={this.setOnClick("Look Away!")} animator={{activate}} hover>{ cornerDiv }</FrameCorners>
                    <FrameCorners onClick={this.setOnClick("Nobody looks at this!")} animator={{activate}} hover>{ cornerDiv }</FrameCorners>
                    <FrameCorners onClick={this.setOnClick("Hey beautiful!")} animator={{activate}} hover>{ cornerDiv }</FrameCorners>
                    <FrameCorners onClick={this.setOnClick("What's the meaning of this?!")} animator={{activate}} hover>{ cornerDiv }</FrameCorners>
                  </nav>
                  <Text as='h1'>Account: { this.state.message }</Text>
                  <section>
                    <Figure src={IMAGE_URL} alt='A nebula'>
                      <Button FrameComponent={FrameHexagon} animator={{ activate }} onClick={() => this.setAssimilate()}>
                        <Text>Assimilate</Text>
                      </Button>
                    </Figure>
                  </section>
                  <div dangerouslySetInnerHTML={{ __html: this.state.text }} />
                </Animator>
              </AnimatorGeneralProvider>
            </BleepsProvider>
          </ArwesThemeProvider>
      </div>
    );
  }
}

export default App;