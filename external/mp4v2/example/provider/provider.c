/* This example makes use of the MP4FileProvider API to use custom file
 * input/output routines.
 */

#include <mp4v2/mp4v2.h>
#include <stdio.h>

/*****************************************************************************/

static void* my_open( const char* name, MP4FileMode mode )
{
    const char* om;
    switch( mode ) {
        case FILEMODE_READ:     om = "rb";  break;
        case FILEMODE_MODIFY:   om = "r+b"; break;
        case FILEMODE_CREATE:   om = "w+b"; break;

        case FILEMODE_UNDEFINED:
        default:
            om = "rb";
            break;
    }

    return fopen( name, om );
}

static int my_seek( void* handle, int64_t pos )
{
    return fseeko( (FILE*)handle, pos, SEEK_SET ) != 0;
}

static int my_read( void* handle, void* buffer, int64_t size, int64_t* nin, int64_t maxChunkSize )
{
    if( fread( buffer, size, 1, (FILE*)handle ) != 1)
        return 1;
    *nin = size;
    return 0;
}

static int my_write( void* handle, const void* buffer, int64_t size, int64_t* nout, int64_t maxChunkSize )
{
    if( fwrite( buffer, size, 1, (FILE*)handle ) != 1)
        return 1;
    *nout = size;
    return 0;
}

static int my_close( void* handle )
{
    return fclose( (FILE*)handle ) != 0;
}

/*****************************************************************************/

int main( int argc, char** argv )
{
    if( argc != 2 ) {
        printf( "usage: %s file.mp4\n", argv[0] );
        return 1;
    }

    /* populate data structure with custom functions.
     * safe to put on stack as it will be immediately copied internally.
     */
    MP4FileProvider provider;

    provider.open  = my_open;
    provider.seek  = my_seek;
    provider.read  = my_read;
    provider.write = my_write;
    provider.close = my_close;

    /* open file for read */
    MP4FileHandle file = MP4ReadProvider( argv[1], 0, &provider );
    if( file == MP4_INVALID_FILE_HANDLE ) {
        printf( "MP4Read failed\n" );
        return 1;
    }

    /* dump file contents */
    if( !MP4Dump( file, stdout, 0 ))
        printf( "MP4Dump failed\n" );

    /* cleanup and close */
    MP4Close( file );

    return 0;
}
