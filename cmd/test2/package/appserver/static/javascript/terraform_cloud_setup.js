"use strict";

// use in production for not logging
// comment out in dev mode if you like
console = {
    log: function() { },
    debug: function() { },
    group: function() { },
    groupEnd: function() { }
};

var app_name = "terraform_cloud";

// Require the setup application
require.config({
    paths: {
        // $SPLUNK_HOME/etc/apps/SPLUNK_APP_NAME/appserver/static/javascript/views/setup
        TFCSetup: "../app/" + app_name + "/javascript/views/setup",
        text: "../app/" + app_name + "/contrib/text",
    },
});

require([
    // Splunk Web Framework Provided files
    "backbone", // From the SplunkJS stack
    "jquery", // From the SplunkJS stack
    // Custom files
    "TFCSetup",
], function(Backbone, jquery, TFCSetup) {
    if (window.frameElement !== null) {
        alert("You need to setup Terraform Cloud for Splunk first.");
        top.location.href = self.location.href;
    }
    var setup_view = new TFCSetup({
        // Sets the element that will be used for rendering
        el: jquery("#main_container"),
    });

    setup_view.render();
});
