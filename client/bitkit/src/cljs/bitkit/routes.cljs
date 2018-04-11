(ns bitkit.routes
  (:require [bidi.bidi :as bidi]
            [re-frame.core :as re-frame]
            [goog.events :as events])
  (:import [goog.history EventType
                         Html5History]))

(defonce routes [["/" :id] :transaction])

(def match (fnil identity {:handler :index}))

(defonce history
  (doto (Html5History.)
    (events/listen EventType.NAVIGATE
      (fn [event]
        (->> (.-token event)
          (bidi/match-route routes)
          match
          :handler
          (vector :set-route)
          re-frame/dispatch)))
    (.setEnabled true)
    (.setUseFragment false)
    (.setPathPrefix "")))

(defn set-token!
  [value]
  (.setToken history value))
