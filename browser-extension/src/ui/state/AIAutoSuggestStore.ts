import { createStore, defaults, StoreActionApi, createHook, createContainer } from 'react-sweet-state';

import { ENDPOINTS } from '../../constants';
import { base64ImageToBlob, createLogger } from '../../utils';

if (process.env.NODE_ENV === 'development') {
  defaults.devtools = true;
}

const log = createLogger('aiAutoSuggestStore');

interface State {
  visible: boolean;
  highlight: boolean;
  loading: boolean;
  error: string | null;
}

const initialState: State = {
  visible: false,
  highlight: false,
  loading: false,
  error: null,
};

interface SuggestContentParams {
  currentUrl: string;
  currentBookingRef: string;
  onSuccess?: (response: string) => void;
  onError?: (error: Error) => void;
}

const actions = {
  setVisible: (visible: boolean) => ({ setState }: StoreActionApi<State>) => setState({ visible }),
  aiSuggestContent:
    ({ currentUrl, currentBookingRef, onSuccess, onError }: SuggestContentParams) =>
      async ({ setState }: StoreActionApi<State>) => {
        setState({ loading: true, highlight: true, error: null });

        try {
          log('taking screenshot', currentUrl);

          const response = await browser.runtime.sendMessage({ action: 'capture_screenshot' });
          const blob = base64ImageToBlob(response.screenshot);

          setState({ highlight: false });

          log('analysing screen context', currentUrl);

          const formData = new FormData();
          formData.append('screenshot', blob, 'screenshot.png');
          let res = await fetch(`${ENDPOINTS.ANALYSE_USER_CLICK}?currentUrl=${encodeURIComponent(currentUrl)}`, {
            method: 'POST',
            body: formData,
          });
          const screenContext = await res.text();

          log('got screen context', screenContext);

          log('requesting ai auto suggest', currentUrl);

          res = await fetch(`${ENDPOINTS.AI_AUTO_SUGGEST}`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              currentUrl,
              currentBookingRef,
              screenContext,
            }),
          });
          const aiSuggestion = await res.text();

          setState({ loading: false, error: null });
          onSuccess?.(aiSuggestion);
        } catch (err: any) {
          setState({ loading: false, error: err.message || 'Unknown error' });
          onError?.(err);
        } finally {
          setState({ highlight: false });
        }
      },
};

export const AIAutoSuggestContainer = createContainer();

export const AIAutoSuggestStore = createStore({
  name: 'AIAutoSuggestLocalStore',
  containedBy: AIAutoSuggestContainer,
  initialState,
  actions,
});

export const useAIAutoSuggestStore = createHook(AIAutoSuggestStore);
