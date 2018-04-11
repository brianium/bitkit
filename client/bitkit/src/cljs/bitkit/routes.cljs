(ns bitkit.routes
  (:require [bidi.bidi :as bidi]
            [re-frame.core :as re-frame]
            [goog.events :as events])
  (:import [goog.history EventType
                         Html5History]))

(defonce routes [["/" :id] :transaction])

(def match (fnil identity {:handler :index}))

(defonce history (Html5History.))

(defn set-path!
  [value]
  (.setToken history value))

(defn init-path! []
  (set-path! (.. js/window -location -pathname)))

(defn listen! []
  (doto history
    (events/listen EventType.NAVIGATE
      (fn [event]
        (->> (.-token event)
          (bidi/match-route routes)
          match
          :handler
          (vector :set-route)
          re-frame/dispatch)))
    (.setUseFragment false)
    (.setPathPrefix "")
    (.setEnabled true))
  (init-path!))
