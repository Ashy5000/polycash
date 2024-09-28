### std00

**Polycash Standards**

*(Draft)*

*Overview*\
In this standard, std01, we define a system for organizing standards in the Polycash ecosystem. Each standard defines a protocol, interface, or meta-standard (a standard that defines the format in which other standards should be produced) that relates to the Polycash ecosystem. They are identified via an ID that ranges from std00 to stdff, with the last two digits forming a hexadecimal number. Each new standard is assigned the next available ID.

*Markdown formatting*\
Each standard must begin with `### stdxx`, where `stdxx` is the standard's ID. Then, its title is written, bolded by adding `**` on either side of it. The first section after the title is an overview, which begins with a section heading (`*Overview*\`). Then, the body of the overview is written using default formatting. After an empty line of Markdown, the other sections are added. Each begins with their header (`*HEADER*\`, with `HEADER` written here as a placeholder for the actual header name). Then, the body of that section follows, again with default formatting.

*Definitions Section*\
Optionally, there can be a `Definitions` section in the standard. Its purpose is to define terms that are used later on in the standard. Its body is composed of a series of lines, each ending with `\` to ensure line breaks are properly respected. Each line begins with the name of the term in backticks (<code>`</code>), followed by a colon, a space, and then the definition of the term. Each definition should not be a complete sentence, but should end in a period.

*State Allocation*
Each standard has a range of local state addresses which it is permitted to use in its protocol. There is an unlimited number of available addresses for each standard. If the standard is stdXX, the standard can include references for how local state should be used from 0xXX to 0xXXffffffffff...unlimited length...f.
For the case of this standard, std00, the available local state, from 0x00 to 0x00ffffffffff...unlimited length...f, is reserved for general use of state; it is not restricted to a specific protocol.