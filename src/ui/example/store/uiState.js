import {
  observable,
  computed,
  asStructure,
  action
} from 'mobx';

import jquery from 'jquery';
import transport from './transport';

//
// Things you will typically find in UI stores:
//
// Session information
// Information about how far your application has loaded
// Information that will not be stored in the backend
// Information that affects the UI globally
//    Window dimensions
//    Accessibility information
//    Current language
//    Currently active theme
// User interface state as soon as it effects multiple, further unrelated components:
//    Current selection
//    Visibility of toolbars, etc.
//    State of a wizard
//    State of a global overlay

class UIState {
  transport;
  @observable language = 'en_US';

  // asStructure makes sure observer won't be signaled only if the
  // dimensions object changed in a deepEqual manner
  @observable windowDimensions = asStructure({
    width: jquery(window).width(),
    height: jquery(window).height()
  });

  @computed get appIsInSync() {
    return this.transport.pending === 0;
  }

  @action resize = () => {
    this.windowDimensions = this.getWindowDimensions();
  }

  constructor(transport) {
    this.transport = transport;

    jquery(window).resize(this.resize);
  }

  getWindowDimensions = () => ({
    width: jquery(window).width(),
    height: jquery(window).height()
  });
}

const uiState = new UIState(transport);

export default uiState;
