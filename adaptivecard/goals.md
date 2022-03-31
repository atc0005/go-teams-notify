# Adaptive Card goals

## TODO

- add a ValidateFunc field to each custom type that exposes a Validate()
  method.

## send2teams prototype

### Equivalent mention (via botapi package)

1. Create Card with user specified text. [first TextBlock]
1. Add mentions to Card.
1. Add trailer to Card [second TextBlock]
1. Set separator to true.

### Equivalent MessageCard functionality

1. Create new MessageCard
1. Set title
1. Set text
1. Set theme color
1. If user specified target URLs, add a new actionSection
1. StartGroup (Q: does this only pertain to seperator lines?) to actionSection
1. Range over target URLs, build PotentialActions, add PotentialActions to actionSection
1. Add actionSection to MessageCard
1. Create new trailerSection
1. Set Text to config.MessageTrailer(...)
1. StartGroup
1. Add trailerSection to MessageCard
1. Send MessageCard to Teams

## Unordered thoughts / Next Steps

continue filling out Validate() functions
add Validate() function for MSTeams type
indicate that MSTeams type is optional by using pointer type?
review other field types to determine which of those should be pointers

~Implement unexported mention() function, use it from Mention() methods.~
Created exported version instead.

TODO: Create multiple helper functions to create useful intermediate values
from common values.
For example:
https://docs.microsoft.com/en-us/azure/bot-service/bot-builder-howto-add-media-attachments?view=azure-bot-service-4.0&tabs=csharp#send-an-adaptive-card
here a JSON file is read and added as the Content (Card) content for a Message
attachment.

Maybe not worth adding NewXYZFromFile variants now, but perhaps something
similar. Perhaps along the lines of NewMessageFromCard().

Add WithSeparator() method to element types that support it. The method is
responsible for enabling the separating line at the top of the element.

Add method or function for adding multiple Mentions to a message. For example,
send2teams accepts a user-specified slice of ID/Name pairs which are intended
to be prepended or appended to a single "message" (e.g., a single TextBlock).

------------

Scratch notes from 2022-03-29 (before heading home):

When adding a mention to an existing Card, should add a new TextBlock element.
When adding a mention to a Message, should add a new TextBlock element to the
first Card.

When adding an existing Mention to a Card, should use the same behavior UNLESS
provided a pointer to an Element. Perhaps this could be a separate standalone
function that accepts pointers to Card (for msteams object access), Element
(Text field). Pointer to Message shouldn't be needed.
