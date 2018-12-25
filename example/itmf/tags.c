/* This is an example of iTMF Tags convenience API.
 * WARNING: this program will change/destroy many tags in an mp4 file.
 */

#include <mp4v2/mp4v2.h>

int main( int argc, char** argv )
{
    if( argc != 2 ) {
        printf( "usage: %s file.mp4\n", argv[0] );
        return 1;
    }

    /* open file for modification */
    MP4FileHandle file = MP4Modify( argv[1], MP4_DETAILS_ERROR, 0 );
    if( file == MP4_INVALID_FILE_HANDLE ) {
        printf( "MP4Modify failed\n" );
        return 1;
    }

    /* allocate */
    const MP4Tags* tags = MP4TagsAlloc();

    /* fetch data from MP4 file and populate structure */
    MP4TagsFetch( tags, file );

    /***************************************************************************
     * print various tag values
     */

    if( tags->name )
        printf( "name: %s\n", tags->name );

    if( tags->artist )
        printf( "artist: %s\n", tags->artist );

    if( tags->albumArtist )
        printf( "albumArtist: %s\n", tags->albumArtist );

    if( tags->album )
        printf( "album: %s\n", tags->album );

    if( tags->grouping )
        printf( "grouping: %s\n", tags->grouping );

    if( tags->composer )
        printf( "composer: %s\n", tags->composer );

    if( tags->comments )
        printf( "comments: %s\n", tags->comments );

    if( tags->genre )
        printf( "genre: %s\n", tags->genre );

    if( tags->genreType )
        printf( "genreType: %u\n", *tags->genreType );

    if( tags->releaseDate )
        printf( "releaseDate: %s\n", tags->releaseDate );

    if( tags->track )
        printf( "track: index=%u total=%u\n", tags->track->index, tags->track->total );

    if( tags->disk )
        printf( "disk: index=%u total=%u\n", tags->disk->index, tags->disk->total );

    if( tags->tempo )
        printf( "tempo: %u\n", *tags->tempo );

    if( tags->compilation )
        printf( "compilation: %u\n", *tags->compilation );

    if( tags->tvShow )
        printf( "tvShow: %s\n", tags->tvShow );

    if( tags->tvNetwork )
        printf( "tvNetwork: %s\n", tags->tvNetwork );

    if( tags->tvEpisodeID )
        printf( "tvEpisodeID: %s\n", tags->tvEpisodeID );

    if( tags->tvSeason )
        printf( "tvSeason: %u\n", *tags->tvSeason );

    if( tags->tvEpisode )
        printf( "tvEpisode: %u\n", *tags->tvEpisode );

    if( tags->description )
        printf( "description: %s\n", tags->description );

    if( tags->longDescription )
        printf( "longDescription: %s\n", tags->longDescription );

    if( tags->lyrics )
        printf( "lyrics: %s\n", tags->lyrics );

    if( tags->sortName )
        printf( "sortName: %s\n", tags->sortName );

    if( tags->sortArtist )
        printf( "sortArtist: %s\n", tags->sortArtist );

    if( tags->sortAlbumArtist )
        printf( "sortAlbumArtist: %s\n", tags->sortAlbumArtist );

    if( tags->sortAlbum )
        printf( "sortAlbum: %s\n", tags->sortAlbum );

    if( tags->sortComposer )
        printf( "sortComposer: %s\n", tags->sortComposer );

    if( tags->sortTVShow )
        printf( "sortTVShow: %s\n", tags->sortTVShow );

    if( tags->artworkCount ) {
        const MP4TagArtwork* art = tags->artwork; /* artwork != NULL when artworkCount > 0 */
        uint32_t i;
        for( i = 0; i < tags->artworkCount; i++, art++ )
            printf( "art[%d]: type=%d size=%u data=%p\n", i, art->type, art->size, art->data );
    }

    if( tags->copyright )
        printf( "copyright: %s\n", tags->copyright );

    if( tags->encodingTool )
        printf( "encodingTool: %s\n", tags->encodingTool );

    if( tags->encodedBy )
        printf( "encodedBy: %s\n", tags->encodedBy );

    if( tags->purchaseDate )
        printf( "purchaseDate: %s\n", tags->purchaseDate );

    if( tags->podcast )
        printf( "podcast: %u\n", *tags->podcast );

    if( tags->keywords )
        printf( "keywords: %s\n", tags->keywords );

    if( tags->category )
        printf( "category: %s\n", tags->category );

    if( tags->hdVideo )
        printf( "hdVideo: %u\n", *tags->hdVideo );

    if( tags->mediaType )
        printf( "mediaType: %u\n", *tags->mediaType );

    if( tags->contentRating )
        printf( "contentRating: %u\n", *tags->contentRating );

    if( tags->gapless )
        printf( "gapless: %u\n", *tags->gapless );

    if( tags->contentID )
        printf( "contentID: %u\n", *tags->contentID );

    if( tags->artistID )
        printf( "artistID: %u\n", *tags->artistID );

    if( tags->playlistID )
        printf( "playlistID: %llu\n", *tags->playlistID );

    if( tags->genreID )
        printf( "genreID: %u\n", *tags->genreID );

    if( tags->xid )
        printf( "xid: %s\n", tags->xid );

    if( tags->iTunesAccount )
        printf( "iTunesAccount: %s\n", tags->iTunesAccount );

    if( tags->iTunesAccountType )
        printf( "iTunesAccountType: %u\n", *tags->iTunesAccountType );

    if( tags->iTunesCountry )
        printf( "iTunesCountry: %u\n", *tags->iTunesCountry );

    /**************************************************************************
     * modify various tags values
     */

    uint8_t  n8;
    uint16_t n16;
    uint32_t n32;
    uint64_t n64;

    MP4TagTrack dtrack;
    dtrack.index = 1;
    dtrack.total = 1;

    MP4TagDisk ddisk;
    ddisk.index = 1;
    ddisk.total = 1;

    MP4TagsSetName              ( tags, "my name" );
    MP4TagsSetArtist            ( tags, "my artist" );
    MP4TagsSetAlbumArtist       ( tags, "my albumArtist" );
    MP4TagsSetAlbum             ( tags, "my album" );
    MP4TagsSetGrouping          ( tags, "my grouping" );
    MP4TagsSetComposer          ( tags, "my composer" );
    MP4TagsSetComments          ( tags, "my comments" );
    MP4TagsSetGenre             ( tags, "my genre" );
    n16 = 5; /* disco */
    MP4TagsSetGenreType         ( tags, &n16 );
    MP4TagsSetReleaseDate       ( tags, "my releaseDate" );
    MP4TagsSetTrack             ( tags, &dtrack );
    MP4TagsSetDisk              ( tags, &ddisk );
    n16 = 60; /* bpm */
    MP4TagsSetTempo             ( tags, &n16 );
    n8 = 0; /* false */
    MP4TagsSetCompilation       ( tags, &n8 );

    MP4TagsSetTVShow            ( tags, "my tvShow" );
    MP4TagsSetTVNetwork         ( tags, "my tvNetwork" );
    MP4TagsSetTVEpisodeID       ( tags, "my tvEpisodeID" );
    n32 = 0;
    MP4TagsSetTVSeason          ( tags, &n32 );
    n32 = 0;
    MP4TagsSetTVEpisode         ( tags, &n32 );

    MP4TagsSetDescription       ( tags, "my description" );
    MP4TagsSetLongDescription   ( tags, "my longDescription" );
    MP4TagsSetLyrics            ( tags, "my lyrics" );

    MP4TagsSetSortName          ( tags, "my sortName" );
    MP4TagsSetSortArtist        ( tags, "my sortArtist" );
    MP4TagsSetSortAlbumArtist   ( tags, "my sortAlbumArtist" );
    MP4TagsSetSortAlbum         ( tags, "my sortAlbum" );
    MP4TagsSetSortComposer      ( tags, "my sortComposer" );
    MP4TagsSetSortTVShow        ( tags, "my sortTVShow" );

    MP4TagsSetCopyright         ( tags, "my copyright" );
    MP4TagsSetEncodingTool      ( tags, "my encodingTool" );
    MP4TagsSetEncodedBy         ( tags, "my encodedBy" );
    MP4TagsSetPurchaseDate      ( tags, "my purchaseDate" );

    n8 = 0; /* false */
    MP4TagsSetPodcast           ( tags, &n8 );
    MP4TagsSetKeywords          ( tags, "my keywords" );
    MP4TagsSetCategory          ( tags, "my category" );

    n8 = 0; /* false */
    MP4TagsSetHDVideo           ( tags, &n8 ); // false
    n8 = 9; /* movie */
    MP4TagsSetMediaType         ( tags, &n8 ); // movie
    n8 = 0; /* none */
    MP4TagsSetContentRating     ( tags, &n8 ); // none
    n8 = 0; /* false */
    MP4TagsSetGapless           ( tags, &n8 ); // false

    MP4TagsSetITunesAccount     ( tags, "my iTunesAccount" );
    n8 = 0; /* iTunes */
    MP4TagsSetITunesAccountType ( tags, &n8 );
    n32 = 143441; /* USA */
    MP4TagsSetITunesCountry     ( tags, &n32 );
    n32 = 0;
    MP4TagsSetContentID         ( tags, &n32 );
    n32 = 0;
    MP4TagsSetArtistID          ( tags, &n32 );
    n64 = 0;
    MP4TagsSetPlaylistID        ( tags, &n64 );
    n32 = 0;
    MP4TagsSetGenreID           ( tags, &n32 );
    n32 = 0;
    MP4TagsSetComposerID        ( tags, &n32 );
    MP4TagsSetXID               ( tags, "my prefix:my scheme:my identifier" );

    /* push data to mp4 file */
    MP4TagsStore( tags, file );

    /* free memory associated with structure and close */
    MP4TagsFree( tags );
    MP4Close( file );

    return 0;
}
