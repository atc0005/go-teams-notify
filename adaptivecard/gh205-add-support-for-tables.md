# GH-205 | Add support for AdaptiveCard Table type

Add types:

- Table
- TableRow
- TableCell

**NOTES**:

- Table is an `Element`
- Column is *not* an Element, but a ColumnSet *is* an Element
- Table entries are shown directly within the `Body` field of a `Card`
  - <https://adaptivecards.io/explorer/Table.html>
- The AdaptiveCard Explorer includes `Table` as an Element type
  - <https://adaptivecards.io/explorer/AdaptiveCard.html#dedupe-headerbody>

Valid (`Element`) types for `Card.Body`:

- `ActionSet`
- `ColumnSet`
- `Container`
- `FactSet`
- `Image`
- `ImageSet`
- `Input.ChoiceSet`
- `Input.Date`
- `Input.Number`
- `Input.Text`
- `Input.Time`
- `Input.Toggle`
- `Media`
- `RichTextBlock`
- `Table`
- `TextBlock`

## Questions

- Does the TableRow type support the ContainerStyle values?

## Current focus

I last looked at whether I needed to include a Style field for TableRow and
TableCell. Both seem to support ContainerStyle values (enums).

Am I planning on inheriting those values somehow? Explicitly adding a field
for those? Presumably styling is supported. Am I making the Element type
handle this for me?

Is the TableRow an Element? Presumably not. Is the TableCell an Element? Also,
presumably not.

For Element, I'll need to make sure that I assert that a Table (as a whole)
does not support a Style field.
