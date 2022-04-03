# Adaptive Card goals

## TODO

- add a ValidateFunc field to each custom type that exposes a Validate()
  method.

## send2teams prototype

### Current MessageCard functionality

See this doc for info on sections:
<https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference#using-sections>

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

### Adaptive Card functionality

1. Create new Message
1. Create new Card
1. Create new Container [initial Container]
1. Create new TextBlock element, set `style` to `heading` [cfg.MessageTitle]
1. Set TextBlock.Text to user-specified Title text
1. Create new TextBlock element [cfg.MessageText]
1. Set TextBlock.Text to user-specified text
1. Add TextBlock to Container
1. Add Container to Card Body
1. If user specified target URLs, add a new Container [for actions]
1. Enable Separator for new Container [actions]
1. Add new ActionSet
1. Range over target URLs
   1. build Action.OpenUrl value
   1. add value to ActionSet
1. Add ActionSet to Container [actions]
1. Add Container to Card Body [actions]
1. Create new Container [trailer]
1. Enable Separator for new Container [trailer]
1. Create new TextBlock element
1. Set TextBlock.Text to config.MessageTrailer(...)
1. Add TextBlock to Container [trailer]
1. Add Container to Card Body [trailer]
1. Add Card to Message
1. Send Message to Teams

### User mention (via botapi package)

1. Create Card with user specified text. [first TextBlock]
1. Add mentions to Card (which modify the first TextBlock)
1. Add trailer to Card [second TextBlock]
1. Set separator to true.

### User mention (via adaptivecard package)

1. Create Card with user specified text. [first TextBlock]
1. Add mentions to Card (which modify the first TextBlock)
1. Add trailer to Card [second TextBlock]
1. Set separator to true.

## Unordered thoughts / Next Steps

It does not appear that an Adaptive Card supports Titles as a MessageCard
does. Instead, examples that I've seen rely on creating an initial TextBlock
element that serves as a Title:

  "body": [
    {
      "type": "Container",
      "items": [
        {
          "type": "TextBlock",
          "text": "Publish Adaptive Card schema",
          "weight": "bolder",
          "size": "medium"
        },

The formatting serves as making the text appear as a Title.

- indicate that MSTeams type is optional by using pointer type?
- review other field types to determine which of those should be pointers

TODO: Create multiple helper functions to create useful intermediate values
from common values.
For example:
<https://docs.microsoft.com/en-us/azure/bot-service/bot-builder-howto-add-media-attachments?view=azure-bot-service-4.0&tabs=csharp#send-an-adaptive-card>
here a JSON file is read and added as the Content (Card) content for a Message
attachment.

Maybe not worth adding NewXYZFromFile variants now, but perhaps something
similar. Perhaps along the lines of NewMessageFromCard().

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
