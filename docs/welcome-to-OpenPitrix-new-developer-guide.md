
Welcome to OpenPitrix! (New Developer Guide)
============================================

_This document assumes that you know what OpenPitrix does. If you don't,
try the demo at [https://o8x.io/](https://o8x.io/)._

Introduction
------------

Have you ever wanted to contribute to the coolest cloud technology? This
document will help you understand the organization of the OpenPitrix project and
direct you to the best places to get started. By the end of this doc, you'll be
able to pick up issues, write code to fix them, and get your work reviewed and
merged.

If you have questions about the development process, feel free to jump into our
[Slack workspace](http://openpitrix.slack.com/) or join our [mailing
list](https://groups.google.com/forum/#!forum/openpitrix-dev). If you join the
Slack workspace it is recommended to set your Slack display name to your GitHub
account handle.

Special Interest Groups
-----------------------

OpenPitrix developers work in teams called Special Interest Groups (SIGs). At
the time of this writing there are [2 SIGs](sig-list.md).

The developers within each SIG have autonomy and ownership over that SIG's part
of OpenPitrix. SIGs organize themselves by meeting regularly and submitting
markdown design documents to the
[openpitrix/community](https://github.com/openpitrix/community) repository.
Like everything else in OpenPitrix, a SIG is an open, community, effort. Anybody
is welcome to jump into a SIG and begin fixing issues, critiquing design
proposals and reviewing code.

Most people who visit the OpenPitrix repository for the first time are
bewildered by the thousands of [open
issues](https://github.com/openpitrix/openpitrix/issues) in our main repository.
But now that you know about SIGs, it's easy to filter by labels to see what's
going on in a particular SIG. For more information about our issue system, check
out
[issues.md](https://github.com/openpitrix/community/blob/master/contributors/devel/issues.md).

//TODO

Downloading, Building, and Testing OpenPitrix
---------------------------------------------

This guide is non-technical, so it does not cover the technical details of
working OpenPitrix. We have plenty of documentation available under
[github.com/openpitrix/openpitrix/docs/](https://github.com/openpitrix/openpitrix/docs/).
Check out
[development.md](https://github.com/openpitrix/openpitrix/docs/development.md)
for more details.

Pull-Request Process
--------------------

The pull-request process is documented in [pull-requests.md](pull-requests.md).
As described in that document, you must sign the CLA before
OpenPitrix can accept your contribution.


The Release Process and Code Freeze
-----------------------------------

Every so often @o8x-merge-robot will refuse to merge your PR, saying something
about release milestones. This happens when we are in a code freeze for a
release. In order to ensure Openpitrix is stable, we stop merging everything
that's not a bugfix, then focus on making all the release tests pass. This code
freeze usually lasts two weeks and happens once per quarter.

If you're new to OpenPitrix, you won't have to worry about this too much. After
you've contributed for a few months, you will be added as a [community
member](https://github.com/openpitrix/openpitrix/docs/membership.md)
and take ownership of some of the tests. At this point, you'll work with members
of your SIG to review PRs coming into your area and track down issues that occur
in tests.

Thanks for reading!
