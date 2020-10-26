# sam - Splunk App Manager

Currently the code in `sbase/` compiles to a binary used for downloading apps from Splunkbase.  See `sbase/README.md` for usage.

---

Deploying Splunk apps to a standalone, or small Splunk install is pretty much a solved issue.  You can search for apps on Splunkbase, deploy them, and update them in a simpole easy UI.

A more complex environment, one with search head clustering, indexer clustering, distributed search, HEC Heavy forwarders, API Heavy Forwarders, Cluster Managers, Deployers, Searhc Head Deployers and Deployment Servers is significantly harder to ensure that the right parts of an app or addon are in the right place, and that you can update them when needed.

Throw into that the desire to adopt more devops processes with testing of apps in the deployment pipeline for both custom and splunkbase apps, also github/gitlab apps, and things get far more comeplex.  If you want to scale this across mulitple sets of indexers, search heads clusters and heavy forwarders, it gets significantly more complex again.

I'm looking to build a way of building a CICD pipeline that does the following:

- [ ] Allows for an app and an addon to be linked if they are distinctly related.
- [ ] Builds App Inspect into the CICD pipeline to fail builds/deployments for apps that do not pass. - <https://dev.splunk.com/enterprise/docs/developapps/testvalidate/appinspect>
- [ ] Repackages apps for specific server functions to ensure each server only gets the pieces of the app/addon combination it needs
- [ ] <https://dev.splunk.com/enterprise/docs/releaseapps/packagingtoolkit> might be a good tool to use to package things for specific types of server deployments.
- [ ] Tracks installed apps and will show when updates are available on the source (Splunkbase, gitlab, github, s3 etc)
- [ ] Push button download, repackage, testing and deployment of apps.

Likely more things as we progress.  Basically, I would like to try and build a more automated, repeatible, simple process for managing Splunk Apps on more comeplxt environments.

## FYI, this code it likely broken more often that it works, it might inface be broken right now
