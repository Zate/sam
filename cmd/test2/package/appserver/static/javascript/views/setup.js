"use strict";
var app_name = "terraform_cloud";

define(
    ["underscore", "backbone", "jquery", "splunkjs/splunk", "text!../app/" + app_name + "/javascript/templates/setupview.html"],
    function(_, Backbone, jquery, splunk_js_sdk, Template) {
        sdk = splunk_js_sdk;

        var TFCView = Backbone.View.extend({
            // -----------------------------------------------------------------
            // Backbone Functions, These are specific to the Backbone library
            // -----------------------------------------------------------------
            initialize: function initialize() {
                console.debug("Changing document title");
                // For some weird reason, the document title is set to terraform_cloud_setup (the name of the js file)
                // so this changes it to something cooler
                window.document.title = "Terraform Cloud for Splunk Setup";
                Backbone.View.prototype.initialize.apply(this, arguments);
            },

            // sets up backbone event listeners for the buttons on the page
            events: {
                "click .setup_button": "trigger_setup",
                "click .skip_button": "manual_setup"
            },

            // default render function
            render: function() {
                console.debug("render: loading template.");
                // Template = ../templates/setupview.html content
                this.el.innerHTML = _.template(Template, {});
                return this;
            },

            // the user has clicked the manual setup button, so we won't setup the modular input
            manual_setup: async function manual_setup() {
                $("#post_error").hide();
                console.group("Manual Setup");
                console.debug("Skipping setup");
                // let's finish up setup since they wanna do modular input stuff themselves
                var result = await this.complete_setup(splunk_js_sdk, false);
                if (!result) {
                    console.groupEnd();
                    return false;
                } else {
                    // everything is setup, we are gonna redirect them
                    console.debug("Setup worked fine, let's send them home");
                    console.groupEnd();
                    await this.redirect_to_home();
                }
                console.groupEnd();
            },

            // the user has clicked the form submit button so we are going to try to get the
            // modular input setup for them
            trigger_setup: async function trigger_setup() {
                $("#post_error").hide();
                console.group("trigger_setup");

                var form = document.querySelector('form');

                // check whether the form is valid using HTML5 validators
                console.debug("Checking form validity");

                if (form.checkValidity()) {
                    console.debug("The form is valid.");

                    // grab the input variables
                    var input_name = jquery("input[name=input_name]").val().trim()
                    var organization_token = jquery("input[name=organization_token]").val().trim();

                    // try to create the modular input
                    var response = await this.create_modular_input(input_name, organization_token);
                    console.debug("Posted with response: ");
                    console.debug(response);

                    // Splunk returns a 200 here no matter what, so we need to check the response.status
                    if (response.status.toLowerCase() == "error") {
                        // there was some error, let's tell the user
                        console.debug("Error: " + response.msg);
                        $("#post_error_message").html(response.msg);
                        $("#post_error").show();
                        console.groupEnd();
                        return false;
                    }
                    // no error, yay!
                    console.debug("No error");
                    $("#post_error").hide();

                    await this.sleep(200);

                    // let's finish up setup since we made the modular input
                    var result = await this.complete_setup(splunk_js_sdk, true);
                    if (!result) {
                        console.groupEnd();
                        return false;
                    } else {
                        // everything is setup, we are gonna redirect them
                        console.debug("Setup worked fine, let's send them home");
                        console.groupEnd();
                        await this.redirect_to_home();
                    }
                }
                console.groupEnd();
            },

            // this is a method that calls to the internal, unpublished manager endpoint (the same one
            // that Splunk uses when you go to Settings->Data->Add Data Input)
            //
            // Splunk may not like this.
            create_modular_input: async function create_modular_input(input_name, organization_token) {
                console.group("create_modular_input")
                console.debug("Input name         = " + input_name);
                console.debug("Organization token = " + organization_token);
                console.groupEnd();

                var data = new FormData();

                data.append("name", input_name);
                data.append("organization_token", organization_token);
                data.append("api_url", "https://app.terraform.io");

                var result = await jQuery.ajax({
                    url: '/en-US/manager/search/data/inputs/terraform_cloud/_new',
                    data: data,
                    async: true,
                    cache: false,
                    contentType: false,
                    processData: false,
                    method: 'POST',
                });

                return JSON.parse(result);
            },

            // this method completes the setup by updating the app.conf file for terraform_cloud
            // setting is_configured to true and then it reloads the application programmatically
            // so that the app.conf is read and it gets out of the setup loop
            complete_setup: async function complete_setup(splunk_js_sdk) {
                console.group("complete_setup");

                var application_name_space = {
                    owner: "nobody",
                    app: app_name,
                    sharing: "app",
                };

                console.debug("application_name_space = " + JSON.stringify(application_name_space))

                var http = new splunk_js_sdk.SplunkWebHttp();

                var service = new splunk_js_sdk.Service(
                    http,
                    application_name_space,
                );

                await this.reload_app(service);
                var result = await this.update_config_file(service);

                if (result) {
                    // this means updating the configuration file was successful, so we reload the app.
                    console.debug("Updating configuration successful");
                    await this.reload_app(service);
                    console.groupEnd();
                    return true;
                } else {
                    // there was some error updating the configuration file, show the error.
                    console.debug("There was an error updating the config file");
                    $("#post_error_message").html("There was an error updating the application configuration. Please check logs.");
                    $("#post_error").show();
                    console.groupEnd();
                    return false;
                }
            },

            sleep: async function sleep(ms) {
                return new Promise(resolve => setTimeout(resolve, ms));
            },
            // this method programmatically calls the Splunk API to update the [install] stanza to reflect
            // that the configuration of the application is complete.  there are a number of checks that
            // happen that are probably unnecessary, but doing it so there are no unexpected errors
            update_config_file: async function update_config_file(service) {
                var configuration_file_name = "app";
                var stanza_name = "install";
                var properties_to_update = {
                    is_configured: "true",
                };

                console.group("update_config_file");
                console.debug("configuration_file_name = " + configuration_file_name);
                console.debug("stanza_name = " + stanza_name);
                console.debug("properties_to_update = " + properties_to_update);

                var config_file_exists = false;
                var stanza_exists = false;

                // this loads up the configuration files for the application
                var configurations = service.configurations({});
                await configurations.fetch();

                var configs = configurations.list();

                // this is a check to make sure that the app.conf config file actually exists.
                // but it totally should and if it doesn't we should be worried
                for (var index = 0; index < configs.length; index++) {
                    var name = configs[index].name;
                    if (name == configuration_file_name) {
                        console.debug("Found " + name);
                        config_file_exists = true;
                    }
                }
                console.debug("config_file_exists = " + config_file_exists);

                if (!config_file_exists) {
                    // this is bad and we should be unhappy
                    // for now, we just return false
                    console.groupEnd();
                    return false;
                }

                // so we know that the app.conf file exists, so we're gonna actually get it
                // (or at least a reference to it)
                var config_file = configurations.item("app", {});
                await config_file.fetch();

                // cool, we have the config file
                console.debug("config_file");
                console.debug(config_file);

                // this will give us all of the stanzas in the configuration file and we'll
                // check that [install] actually exists.
                // it should exist.
                var stanzas = config_file.list();
                for (var index = 0; index < stanzas.length; index++) {
                    var name = stanzas[index].name;
                    if (name === stanza_name) {
                        console.debug("Found " + name);
                        stanza_exists = true;
                    }
                }
                console.debug("stanza_exists = " + stanza_exists);

                // this means that the install stanza doesn't exist and we should be unhappy
                if (!stanza_exists) {
                    console.groupEnd();
                    return false;
                }

                // the install stanza exists so we're gonna get it (or a reference to it) so that
                // we can update it.
                var stanza = config_file.item(stanza_name, {});
                await stanza.fetch();

                // this is the call to update the stanza.  the error response function is just for
                // logging - we get a response back and put it in stanza_result and log that.
                var stanza_result = await stanza.update(
                    properties_to_update,
                    function(error_response, entity) {
                        if (error_response != null) {
                            console.debug("There was an error updating the stanza");
                            console.debug(error_response);
                        } else {
                            console.debug("Successfully updated stanza.");
                        }
                    },
                );

                // let's log this out for fun
                console.debug("stanza_result");
                console.debug(stanza_result);

                // we are going to reload the stanza from the server and see if our update (0 -> 1)
                // happened or not.  all of the above checks were to get to this meat and potatoes
                await stanza.fetch();
                if (stanza.state().content.is_configured !== "1") {
                    // oh no - we weren't able to update the install stanza setting.
                    console.groupEnd();
                    return false;
                }

                // we are cool, let's move on
                console.groupEnd();
                return true;
            },

            // this method reloads the application so that it reloads the application configuration
            reload_app: async function reload_app(service) {
                console.debug("Reloading application");

                var apps = service.apps();
                await apps.fetch();

                var current_app = apps.item(app_name);
                current_app.reload();
                current_app.reload();
            },

            // this will send the user to the main page of the application
            redirect_to_home: async function redirect_to_home() {
                var redirect_url = "/app/" + app_name;
                console.debug("Redirecting to " + redirect_url);
                window.setTimeout(() => {
                    console.debug("Here comes the redirect");
                    window.location.href = redirect_url;
                }, 2500);
            }
        });

        return TFCView;
    }, // End of require asynchronous module definition function
); // End of require statement
