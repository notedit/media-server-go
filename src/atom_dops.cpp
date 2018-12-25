/*
 * The contents of this file are subject to the Mozilla Public
 * License Version 1.1 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of
 * the License at http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS
 * IS" basis, WITHOUT WARRANTY OF ANY KIND, either express or
 * implied. See the License for the specific language governing
 * rights and limitations under the License.
 *
 * The Original Code is MPEG4IP.
 *
 * The Initial Developer of the Original Code is Cisco Systems Inc.
 * Portions created by Cisco Systems Inc. are
 * Copyright (C) Cisco Systems Inc. 2004.  All Rights Reserved.
 *
 * Contributor(s):
 *      Bill May wmay@cisco.com
 */

#include "src/impl.h"

namespace mp4v2 {
namespace impl {

///////////////////////////////////////////////////////////////////////////////

MP4DOpsAtom::MP4DOpsAtom(MP4File &file)
        : MP4Atom(file, "dOps")
{
    
    AddProperty( new MP4Integer8Property(*this,"version")); /* 0 */
    AddProperty( new MP4Integer8Property(*this,"outputChannelCount")); /* 1 */
    AddProperty( new MP4Integer16Property(*this,"preSkip")); /* 2 */
    AddProperty( new MP4Integer32Property(*this,"inputSampleRate")); /* 3 */
    AddProperty( new MP4Integer16Property(*this,"outputGain")); /* 4 */
    AddProperty( new MP4Integer8Property(*this,"channelMappingFamily")); /* 5 */
}

void MP4DOpsAtom::Generate()
{
    MP4Atom::Generate();
    ((MP4Integer8Property*)m_pProperties[0])->SetValue(0);
    ((MP4Integer8Property*)m_pProperties[1])->SetValue(2);
    ((MP4Integer16Property*)m_pProperties[2])->SetValue(0);
    ((MP4Integer32Property*)m_pProperties[3])->SetValue(48000);
    ((MP4Integer32Property*)m_pProperties[4])->SetValue(0);
    ((MP4Integer32Property*)m_pProperties[5])->SetValue(0);
}
///////////////////////////////////////////////////////////////////////////////

}
} // namespace mp4v2::impl
